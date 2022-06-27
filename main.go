package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"main/models"
	"main/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func AddProductRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		buf     bytes.Buffer
		modelIn models.ProductIn
		images  []models.ImageReceptProduct
	)
	io.Copy(&buf, r.Body)
	json.Unmarshal(buf.Bytes(), &modelIn)
	fmt.Printf("Add new product \"%s\"\n", modelIn.Name)
	product := models.Product{Name: modelIn.Name}
	err := db.Create(&product).Error
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	for _, image := range modelIn.Images {
		var image = models.ImageReceptProduct{
			Image:     image,
			ProductId: product.ID,
		}
		db.Create(&image)
		images = append(images, image)
	}
	result, _ := json.Marshal(models.BaseResponse{
		Result: models.ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Images: images,
		},
	})
	w.Write(result)
}

func GetImageRouterHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	r.ParseForm()
	fileName := r.Form.Get("file")
	file, err := os.Open("./images/" + fileName)
	if err != nil || buf.Len() != 0 {
		fmt.Fprint(w, "No file")
		return
	}
	io.Copy(&buf, file)
	w.Write(buf.Bytes())
}

func AddReceptRouterHadler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		buf      bytes.Buffer
		receptIn models.ReceptIn
		author   []models.User
	)
	io.Copy(&buf, r.Body)
	json.Unmarshal(buf.Bytes(), &receptIn)
	fmt.Printf("Add new recept: \"%s\":\"%s\"\n", receptIn.Name, receptIn.Author)

	db.Where("g_id = ?", receptIn.Author).Find(&author)

	if len(author) == 0 {
		result, _ := json.Marshal(models.BaseResponse{
			Error: "No such user",
		})
		w.Write(result)
		return
	}

	recept := models.Recept{
		Name:        receptIn.Name,
		Description: receptIn.Description,
		Author:      author[0].ID,
	}
	err := db.Create(&recept).Error
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	for _, productID := range receptIn.Products {
		db.Create(&models.ReceptProduct{
			ProducId: productID,
			ReceptId: recept.ID,
		})
	}
	for _, image := range receptIn.Images {
		db.Create(&models.ImageReceptProduct{
			Image:    image,
			ReceptId: recept.ID,
		})
	}
	result, _ := json.Marshal(models.BaseResponse{Result: recept.ID})
	w.Write(result)
}

func UploadImageRouterHadler(w http.ResponseWriter, r *http.Request) {
	imageID, err := utils.SaveFile(r)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	result, _ := json.Marshal(models.ImageResponse{ImageID: "http://nextrun.keenetic.pro:8080/get_image?file=" + imageID})
	w.Write(result)
}

// Выдача рецептов по 20 едениц
func QueryReceptRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		recepts       []models.Recept
		resultRecepts []models.ReceptResponse
	)
	fmt.Println("QueryReceptRouterHandler")
	page, err := utils.ToInt(r.Form.Get("page"))
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	offset := 20 * page
	db.Limit(offset + 20).Offset(offset).Find(&recepts)
	for _, recept := range recepts {
		var (
			images       []models.ImageReceptProduct
			products_ids []models.ReceptProduct
			products     []models.ProductResponse
			user         models.User
		)

		db.Where("recept_id = ?", recept.ID).Find(&images)
		db.Where("recept_id = ?", recept.ID).Find(&products_ids)
		db.Find(&user, recept.Author)

		for _, product_id := range products_ids {
			var (
				product        models.Product
				product_images []models.ImageReceptProduct
			)
			db.Find(&product, product_id.ProducId)
			db.Where("product_id = ?", product.ID).Find(&product_images)
			products = append(products, models.ProductResponse{
				ID:     product.ID,
				Name:   product.Name,
				Images: product_images,
			})
		}

		resultRecepts = append(resultRecepts, models.ReceptResponse{
			ID:       recept.ID,
			Author:   user,
			Name:     recept.Name,
			Images:   images,
			Products: products,
		})
	}
	result, _ := json.Marshal(models.BaseResponse{Result: resultRecepts})
	w.Write(result)
}

// Выдача продуктов по 20 едениц
func QueryProductRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		products       []models.Product
		resultProducts []models.ProductResponse
	)
	fmt.Println("QueryProductRouterHandler")
	page, err := utils.ToInt(r.Form.Get("page"))
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	offset := page * 20
	err1 := db.Limit(offset + 20).Offset(offset).Find(&products).Error
	if err1 != nil {
		fmt.Fprint(w, err1.Error())
		return
	}
	for _, product := range products {
		var images []models.ImageReceptProduct
		err2 := db.Where("product_id = ?", product.ID).Find(&images).Error
		if err2 != nil {
			fmt.Fprint(w, err2.Error())
			return
		}
		resultProducts = append(resultProducts, models.ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Images: images,
		})
	}
	resutl, _ := json.Marshal(models.BaseResponse{Result: resultProducts})
	w.Write(resutl)
}

func GetProductsRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		connects []models.ReceptProduct
		products []models.ProductResponse
	)
	fmt.Println("GetProductsRouterHandler")
	receptID, err := utils.ToInt(r.Form.Get("receptID"))
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	db.Find(&connects, "recept_id = ?", receptID)
	for _, productID := range connects {
		var (
			product models.Product
			images  []models.ImageReceptProduct
		)
		db.Find(&product, productID.ProducId)
		db.Where("product_id = ?", productID.ProducId).Find(&images)
		products = append(products, models.ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Images: images,
		})
	}
	result, _ := json.Marshal(models.BaseResponse{Result: products})
	w.Write(result)
}

func GetProductByNameRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var (
		name           = r.Form.Get("search")
		products       []models.Product
		productsResult []models.ProductResponse
	)
	fmt.Printf("Search: %s\n", name)
	db.Where("name LIKE ?", "%"+name+"%").Find(&products)
	for _, product := range products {
		var images []models.ImageReceptProduct
		db.Where("product_id = ?", product.ID).Find(&images)
		productsResult = append(productsResult, models.ProductResponse{
			ID:     product.ID,
			Name:   product.Name,
			Images: images,
		})
	}
	result, _ := json.Marshal(models.BaseResponse{Result: productsResult})
	w.Write(result)
}

func PoliticsRouterHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	file, err := os.Open("./politics/politics.html")
	if err != nil {
		fmt.Fprint(w, "No file")
		return
	}
	io.Copy(&buf, file)
	w.Write(buf.Bytes())
}

func AddNewUserRouterHandler(w http.ResponseWriter, r *http.Request) {
	var (
		buf    bytes.Buffer
		user   models.UserIn
		dbUser []models.User
	)
	io.Copy(&buf, r.Body)
	err := json.Unmarshal(buf.Bytes(), &user)
	db.Where("g_id = ?", user.GID).Find(&dbUser)
	if len(dbUser) != 0 {
		result, _ := json.Marshal(models.BaseResponse{
			Result: models.UserResponse{
				User: dbUser[0],
				New:  false,
			},
		})
		w.Write(result)
		return
	}
	if err != nil {
		result, _ := json.Marshal(models.BaseResponse{
			Error: err.Error(),
		})
		w.Write(result)
		return
	}
	userResult := models.User{
		Name:  user.Name,
		Image: user.Image,
		GID:   user.GID,
	}
	err1 := db.Create(&userResult).Error
	if err1 != nil {
		result, _ := json.Marshal(models.BaseResponse{
			Error: err1.Error(),
		})
		w.Write(result)
	}
	result, _ := json.Marshal(models.BaseResponse{
		Result: models.UserResponse{
			User: userResult,
			New:  true,
		},
	})
	w.Write(result)
}

func UpdateUserRouterHanler(w http.ResponseWriter, r *http.Request) {
	var (
		buf    bytes.Buffer
		user   models.UserIn
		dbUser []models.User
	)
	io.Copy(&buf, r.Body)
	json.Unmarshal(buf.Bytes(), &user)
	db.Where("g_id = ?", user.GID).Find(&dbUser)
	if len(dbUser) == 0 {
		result, _ := json.Marshal(models.BaseResponse{
			Error: "No such user",
		})
		w.Write(result)
		return
	}
	var resultUser = models.User{
		Name:  user.Name,
		Image: user.Image,
		GID:   dbUser[0].GID,
	}
	resultUser.ID = dbUser[0].ID
	db.Save(&resultUser)
	db.Where("g_id = ?", user.GID).Find(&dbUser)
	result, _ := json.Marshal(models.BaseResponse{
		Result: dbUser[0],
	})
	w.Write(result)
}

func main() {
	fmt.Println("Server started!")
	dbt, dbErr := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = dbt

	if dbErr != nil {
		panic("No database!")
	}

	db.AutoMigrate(
		&models.Product{},
		&models.Recept{},
		&models.ImageReceptProduct{},
		&models.ReceptProduct{},
		&models.User{},
	)

	http.HandleFunc("/add_product", AddProductRouterHandler)          // Добавление продукта в бд
	http.HandleFunc("/add_recept", AddReceptRouterHadler)             // Добавление рецепта в бд
	http.HandleFunc("/add_image", UploadImageRouterHadler)            // Добавление изображения
	http.HandleFunc("/get_image", GetImageRouterHandler)              // Получение изображения
	http.HandleFunc("/query_recept", QueryReceptRouterHandler)        // Получение рецептов
	http.HandleFunc("/query_product", QueryProductRouterHandler)      // Получение продуктов
	http.HandleFunc("/search_product", GetProductByNameRouterHandler) // Получение продуктов по названию
	http.HandleFunc("/log_in", AddNewUserRouterHandler)               // Добавление нового пользователя
	http.HandleFunc("/update_user", UpdateUserRouterHanler)           // Обновление пользователя

	http.HandleFunc("/politics", PoliticsRouterHandler) // Политика кондфициальности

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
