package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mongdbs/database"
	"mongdbs/model"
	"mongdbs/resolvers"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Body struct {
	Type string   `json:"type" form:"type"`
	Id   string   `json:"id" form:"id"`
	Name []string `json:"name" form:"name"`
}

var (
	IMAGE_FOLDER = "/var/data/images" // 图片存储的固定路径
	temp         = "ret2-image-temp-folder"
)

func main() {
	database.InitDB_docker()
	// database.InitDB()
	r := gin.Default()

	// 定义各个路由和对应的处理函数
	r.POST("/api/knowledge", createKnowledgeHandler)
	r.PUT("/api/knowledge/:id", updateKnowledgeHandler)
	r.DELETE("/api/knowledge/:id", deleteKnowledgeHandler)
	r.GET("/api/knowledge/type", searchByKnowledgeTypeHandler)
	r.GET("/api/knowledge/tactics", mitreByTacticsIDHandler)
	r.GET("/api/knowledge/techniques", mitreByTechniquesIDHandler)
	r.GET("/api/knowledge/subtechniques", mitreBySubTechniquesIDHandler)
	r.GET("/api/knowledge/search", searchHandler)

	r.GET("/api/knowledge/title", searchByTitleHandler) // 空格需要被替换成为%20
	r.GET("/api/knowledge/tags", searchByTagsWithTypeHandler)
	r.GET("/api/knowledge/content", searchByContentHandler)
	r.GET("/api/knowledge/keyword", searchByKeywordHandler)
	r.GET("/api/knowledge/id", searchByIDHandler) // 新添加的通过ID查询路由
	r.POST("/api/knowledge/batchEdit", batchEditKnowledgeTypeHandler)

	// 图片处理相关路由
	r.POST("/api/images", UploadImageHandler)
	r.POST("/api/path", UploadImagePath)
	r.GET("/api/images/:type/:id/:filename", DownloadFileHandler)
	r.DELETE("/api/images/:type/:id/:filename", DeleteImageHandler)

	log.Println("Server is running on port 8085")
	r.Run(":8085")
	// r.POST("/api/knowledge/search", searchHandler)
}

func UploadImageHandler(c *gin.Context) {
	fileFolder := IMAGE_FOLDER

	file, handler, err := c.Request.FormFile("images")
	if err != nil {
		log.Printf("Error retrieving form file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 验证文件类型
	imageType := strings.ToLower(handler.Header.Get("Content-Type"))
	log.Printf("imageType: %v\n", imageType)
	switch imageType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png", "image/tiff", "image/eps", "image/svg", "image/pdf", "image/bmp":
		log.Println("上传的图片类型:", imageType)
	default:
		log.Println("Invalid image type:", imageType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type"})
		return
	}

	var body Body
	if err := c.Bind(&body); err != nil {
		log.Printf("Error binding request payload: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	filename := handler.Filename
	ext := filepath.Ext(filename)
	if ext == "" {
		log.Println("Invalid file extension")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	uuId := uuid.New().String()
	currentTime := time.Now()
	newFileName := fmt.Sprintf("%04d-%02d-%02d_%02d%02d%02d_%s%s",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second(), uuId, ext)

	log.Printf("Generated new file name: %s\n", newFileName)

	subPath := ""
	if body.Id != "" {
		subPath = filepath.Join(body.Type, body.Id, newFileName)
	} else {
		subPath = filepath.Join(body.Type, temp, newFileName)
	}

	log.Printf("Sub path for file: %s\n", subPath)

	dstPath := filepath.Join(fileFolder, subPath)
	log.Printf("Destination path: %s\n", dstPath)

	err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating directory"})
		return
	}
	log.Printf("Directory created successfully: %s\n", filepath.Dir(dstPath))

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Error creating file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file"})
		return
	}
	defer dst.Close()
	log.Printf("File created successfully: %s\n", dstPath)

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error copying file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying file"})
		return
	}
	log.Println("File copied successfully")

	// 返回图片的访问URL，修改为你的服务域名或IP地址
	c.JSON(http.StatusOK, gin.H{
		"message": "图片上传成功",
		"image":   fmt.Sprintf("http://192.168.31.56:8085/api/images/%s", subPath),
	})
}

// UploadImagePath 接收id参数，将暂存的图片移动到/images/knowledgeType/id
func UploadImagePath(c *gin.Context) {
	log.Println("UploadImagePath called")
	log.Printf("Request data: %v\n", c.Request)
	fileFolder := IMAGE_FOLDER

	var body Body
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	mkdirPath := filepath.Join(fileFolder, body.Type, body.Id)
	err := os.MkdirAll(mkdirPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	for _, filename := range body.Name {
		srcPath := filepath.Join(fileFolder, body.Type, "temp", filename)
		dstPath := filepath.Join(fileFolder, body.Type, body.Id, filename)

		err := copyFile(srcPath, dstPath)
		if err != nil {
			log.Println("Error copying file: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying file"})
			return
		}

		err = os.RemoveAll(srcPath)
		if err != nil {
			log.Println("Error removing temp file: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing temp file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": true,
	})
}

// DeleteImageHandler 处理图片删除
func DeleteImageHandler(c *gin.Context) {
	fileFolder := IMAGE_FOLDER

	filename := c.Param("filename")
	ftype := c.Param("type")
	name := c.Param("id")

	if filename == "" || ftype == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	subPath := filepath.Join(ftype, name, filename)
	dstPath := filepath.Join(fileFolder, subPath)

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if err := os.Remove(dstPath); err != nil {
		log.Println("Error deleting file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file", "detail": err.Error()})
		return
	}

	dirPath := filepath.Join(fileFolder, ftype, name)
	if err := os.Remove(dirPath); err != nil {
		// 忽略文件夹删除错误，因为文件夹可能非空
		log.Println("Error deleting directory (ignored): ", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// DownloadFileHandler 处理图片下载
func DownloadFileHandler(c *gin.Context) {
	fileFolder := IMAGE_FOLDER

	filename := c.Param("filename")
	ftype := c.Param("type")
	name := c.Param("id")

	if filename == "" || ftype == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	subPath := filepath.Join(ftype, name, filename)
	dstPath := filepath.Join(fileFolder, subPath)

	_, err := os.Stat(dstPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	file, err := os.Open(dstPath)
	if err != nil {
		log.Println("Error opening file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(dstPath)
	ext = strings.TrimPrefix(ext, ".")
	log.Println("ext: ", ext)

	// c.Header("Content-Type", "image/"+ext+"")
	contentType := "application/octet-stream" // 默认的Content-Type
	switch ext {
	case "svg":
		contentType = "image/svg+xml"
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "gif":
		contentType = "image/gif"
	case "bmp":
		contentType = "image/bmp"
	case "tiff":
		contentType = "image/tiff"
	case "eps":
		contentType = "application/postscript"
	case "pdf":
		contentType = "application/pdf"
	}

	c.Header("Content-Type", contentType)
	http.ServeContent(c.Writer, c.Request, filename, time.Time{}, file)
}

// copyFile is a utility function to copy files from source to destination
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func batchEditKnowledgeTypeHandler(c *gin.Context) {
	var req struct {
		IDList   []string `json:"idList"`
		PrevType string   `json:"prevType"`
		RepType  string   `json:"repType"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	status, err := resolver.Mutation().BatchEditKnowledgeType(ctx, req.IDList, req.PrevType, req.RepType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// 处理函数 测试没问题
func createKnowledgeHandler(c *gin.Context) {
	var knowledge model.NewKnowledge
	if err := c.BindJSON(&knowledge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	createdKnowledge, err := resolver.Mutation().CreateKnowledge(ctx, knowledge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("success")
	log.Println("Successfully created knowledge:", createdKnowledge)
	c.JSON(http.StatusOK, createdKnowledge)
}

// 测试没问题
func updateKnowledgeHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("id: %v\n", id)
	var knowledge model.NewKnowledge
	if err := c.BindJSON(&knowledge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	updatedKnowledge, err := resolver.Mutation().UpdateKnowledge(ctx, id, knowledge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedKnowledge)
}

// 测试没问题
func deleteKnowledgeHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("id: %v\n", id)
	ctx := context.Background()
	resolver := resolvers.Resolver{}
	deletionStatus, err := resolver.Mutation().DeleteKnowledge(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deletionStatus)
}

// 测试没问题
func searchByKnowledgeTypeHandler(c *gin.Context) {
	typeArgs := c.QueryArray("type")
	numsStr := c.DefaultQuery("nums", "0")
	fmt.Printf("typeArgs: %v\n", typeArgs)
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchByKnowledgeType(ctx, typeArgs, nums)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// 测试没问题
func mitreByTacticsIDHandler(c *gin.Context) {
	tacticsID := c.QueryArray("tacticsId") // 得到id
	fmt.Printf("tacticsID: %v\n", tacticsID)
	if len(tacticsID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tacticsId must be provided"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().MitreByTacticsID(ctx, tacticsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// 测试没问题
func mitreByTechniquesIDHandler(c *gin.Context) {
	// techniquesID := c.QueryArray("techniquesID")
	// curl -X GET "http://localhost:8085/api/knowledge/techniques?techniquesId=tech1" 注意请求参数要完全要一致

	// if len(techniquesID) == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "techniquesID must be provided"})
	// 	return
	// }
	techniquesID := c.QueryArray("TechniquesId")
	fmt.Println("techniquesID: ", techniquesID)
	if len(techniquesID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TechniquesId must be provided"})
		return
	}
	// 这个ID有问题有一些不可见字符
	cleanedTechniquesID := make([]string, len(techniquesID))

	for i, id := range techniquesID {
		cleanedTechniquesID[i] = cleanString(id)
		fmt.Printf("Cleaned TechniquesId: '%s'\n", cleanedTechniquesID[i])
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().MitreByTechniquesID(ctx, cleanedTechniquesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func mitreBySubTechniquesIDHandler(c *gin.Context) {
	subTechniquesID := c.QueryArray("SubTechniquesId")

	if len(subTechniquesID) == 0 {
		singleID := c.Query("SubTechniquesId")
		if singleID != "" {
			subTechniquesID = []string{singleID}
		}
	}

	fmt.Printf("Processed subTechniquesID: %v\n", subTechniquesID)

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().MitreBySubTechniquesID(ctx, subTechniquesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// 组合查询方法 目前没问题
func searchHandler(c *gin.Context) {
	var where model.KnowledgeFilter
	where.Tags = c.QueryArray("tags") // 预过滤在这些匹配项中寻找 keyword相匹配的关键项目

	// where.Tags = c.QueryArray("tags")
	where.Confidentiality = c.Query("confidentiality")
	where.KnowledgeType = c.QueryArray("knowledgeType")

	// Add additional fields
	where.Title = c.Query("title")
	where.Abstract = c.Query("abstract")
	where.Content = c.Query("content")
	where.Recommendations = c.Query("recommendations")
	where.Solution = c.Query("solution")
	where.TacticsID = c.QueryArray("tacticsId")
	where.TechniquesID = c.QueryArray("techniquesId")
	where.SubTechniquesID = c.QueryArray("subTechniquesId")

	authors := c.QueryArray("author")
	keyword := c.QueryArray("keyword")
	nodedict := c.Query("nodedict")
	fmt.Printf("nodedict: %v\n", nodedict)
	numsStr := c.DefaultQuery("nums", "0")
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}
	log.Printf("Searching with tags: %v, keyword: %v, nums: %d", where.Tags, keyword, nums)

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().Search(ctx, &where, keyword, authors, nums, nodedict)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// 使用%20 来代替里面出现的空格curl -X GET "http://localhost:8085/api/knowledge/title?title=Knowledge%201&nums=5"

func searchByTitleHandler(c *gin.Context) {
	title := c.Query("title")
	fmt.Printf("title: %v\n", title)
	numsStr := c.DefaultQuery("nums", "0")
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchByTitle(ctx, title, nums)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// curl -X GET "http://localhost:8085/api/knowledge/tags?type=type1&type=type2&tags=tag1&tags=tag2&nums=5" 测试至少没问题
func searchByTagsWithTypeHandler(c *gin.Context) {
	typeArg := c.QueryArray("type")
	tags := c.QueryArray("tags")
	numsStr := c.DefaultQuery("nums", "0")
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchByTagsWithType(ctx, typeArg, tags, nums)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// curl -X GET "http://localhost:8085/api/knowledge/content?type=type1&keyword=content&nums=2" keyword是要搜索的内容在content中
func searchByContentHandler(c *gin.Context) {
	typeArg := c.QueryArray("type")
	keyword := c.Query("keyword")
	numsStr := c.DefaultQuery("nums", "0")
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchByContent(ctx, typeArg, keyword, nums)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// curl -X GET "http://localhost:8085/api/knowledge/keyword?type=type1&keyword=content&nums=2" 按关键字在某字段搜索
func searchByKeywordHandler(c *gin.Context) {
	typeArg := c.QueryArray("type")
	keyword := c.QueryArray("keyword")
	numsStr := c.DefaultQuery("nums", "0")
	nums, err := strconv.Atoi(numsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid nums parameter"})
		return
	}

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchByKeyword(ctx, typeArg, keyword, nums)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func searchByIDHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be provided"})
		return
	}
	typeArg := c.QueryArray("type")

	ctx := context.Background()
	resolver := resolvers.Resolver{}
	results, err := resolver.Query().SearchById(ctx, typeArg, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func cleanString(input string) string {
	// 去除所有不可见字符
	cleaned := strings.ReplaceAll(input, "\u200B", "")
	cleaned = strings.ReplaceAll(cleaned, "\u200C", "")
	cleaned = strings.ReplaceAll(cleaned, "\u200D", "")
	cleaned = strings.ReplaceAll(cleaned, "\uFEFF", "")
	return cleaned
}
