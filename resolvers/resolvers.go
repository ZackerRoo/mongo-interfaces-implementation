package resolvers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mongdbs/database"
	"mongdbs/model"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Resolver struct{}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

type MutationResolver interface {
	CreateKnowledge(ctx context.Context, input model.NewKnowledge) (*model.Knowledge, error)
	UpdateKnowledge(ctx context.Context, id string, input model.NewKnowledge) (*model.Knowledge, error)
	DeleteKnowledge(ctx context.Context, id string) (*model.DeletionStatus, error)
	BatchEditKnowledgeType(ctx context.Context, idList []string, prevType string, repType string) (*model.DeletionStatus, error)
}

type QueryResolver interface {
	SearchByKnowledgeType(ctx context.Context, typeArg []string, nums int) ([]*model.Knowledge, error)
	MitreByTacticsID(ctx context.Context, tacticsID []string) ([]*model.Knowledge, error)
	MitreByTechniquesID(ctx context.Context, techniquesID []string) ([]*model.Knowledge, error)
	MitreBySubTechniquesID(ctx context.Context, subTechniquesID []string) ([]*model.Knowledge, error)
	Search(ctx context.Context, where *model.KnowledgeFilter, keyword []string, authors []string, nums int, nodedict string) ([]*model.Knowledge, error)
	SearchByTitle(ctx context.Context, title string, nums int) ([]*model.Knowledge, error)
	SearchByTagsWithType(ctx context.Context, typeArg []string, tags []string, nums int) ([]*model.Knowledge, error)
	SearchByContent(ctx context.Context, typeArg []string, keyword string, nums int) ([]*model.Knowledge, error)
	SearchByKeyword(ctx context.Context, typeArg []string, keyword []string, nums int) ([]*model.Knowledge, error)
	SearchById(ctx context.Context, typeArg []string, id string) ([]*model.Knowledge, error)
}

func (r *mutationResolver) BatchEditKnowledgeType(ctx context.Context, idList []string, prevType string, repType string) (*model.DeletionStatus, error) {
	collection := database.GetCollection("knowledge")

	if len(idList) <= 0 {
		return &model.DeletionStatus{
			Success: false,
			Message: "id list 为空",
		}, nil
	}

	var fail []string
	for _, id := range idList {
		// 查找文档
		var knowledge model.Knowledge
		err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&knowledge)
		if err != nil {
			fail = append(fail, id+": 不存在")
			continue
		}

		// 替换knowledgeType
		found := false
		for i, t := range knowledge.KnowledgeType {
			if repType == t {
				continue
			}
			if prevType == t {
				knowledge.KnowledgeType[i] = repType
				found = true
				break
			}
		}

		if !found {
			fail = append(fail, id+": 未找到prevType")
			continue
		}

		// 更新文档 id 字段为_id
		update := bson.M{"$set": bson.M{"knowledgeType": knowledge.KnowledgeType}}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, update)
		if err != nil {
			fail = append(fail, id+": 更新失败，"+err.Error())
			continue
		}

		// 验证更新结果
		var updatedKnowledge model.Knowledge
		err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedKnowledge)
		if err != nil || !reflect.DeepEqual(updatedKnowledge.KnowledgeType, knowledge.KnowledgeType) {
			result := fmt.Sprintf("[%s]", strings.Join(strings.Fields(fmt.Sprint(updatedKnowledge.KnowledgeType)), ", "))
			fail = append(fail, id+": 更新失败，结果为："+result)
		}
	}

	if len(fail) > 0 {
		return &model.DeletionStatus{
			Success: false,
			Message: "部分更新失败：" + strings.Join(fail, "; "),
		}, nil
	}

	return &model.DeletionStatus{
		Success: true,
		Message: "成功",
	}, nil
}

// 一般来说找不到的话不要报错直接返回空比较好
func (r *queryResolver) SearchById(ctx context.Context, typeArg []string, id string) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"_id": id}
	if len(typeArg) > 0 && typeArg[0] != "" {
		filter["knowledgeType"] = bson.M{"$in": typeArg}
	}

	var results []*model.Knowledge
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result model.Knowledge
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		result.Success = true
		results = append(results, &result)
	}

	if len(results) == 0 {
		return results, nil
	}

	return results, nil
}

// CreateKnowledge 实现
// curl -X POST http://localhost:8085/api/knowledge -H "Content-Type: application/json" -d '{
//     "id": "5",
//     "title": "Knowledge 5",
//     "tags": ["tag5", "tag1"],
//     "techniquesId": ["tech5"],
//     "tacticsId": ["tactic5"],
//     "knowledgeType": ["type5"],
//     "content": "This is the content for Knowledge 5."
// }'

// curl -X POST http://localhost:8085/api/knowledge -H "Content-Type: application/json" -d '{
//     "id": "6",
//     "title": "Knowledge 6",
//     "tags": ["tag6", "tag2"],
//     "techniquesId": ["tech6"],
//     "tacticsId": ["tactic6"],
//     "knowledgeType": ["type6"],
//     "content": "This is the content for Knowledge 6."
// }'

// curl -X POST http://localhost:8085/api/knowledge -H "Content-Type: application/json" -d '{
//     "id": "7",
//     "title": "Knowledge 7",
//     "tags": ["tag7", "tag3"],
//     "techniquesId": ["tech7"],
//     "tacticsId": ["tactic7"],
//     "knowledgeType": ["type7"],
//     "content": "This is the content for Knowledge 7."
// }'

func (r *mutationResolver) CreateKnowledge(ctx context.Context, input model.NewKnowledge) (*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	// 将 _id 设置为 input.ID
	if input.ID == "" {
		input.ID = primitive.NewObjectID().Hex()
	}

	doc := model.Knowledge{
		ID:                  input.ID,
		Title:               input.Title,
		Tags:                input.Tags,
		TechniquesID:        input.TechniquesID,
		TacticsID:           input.TacticsID,
		KnowledgeType:       input.KnowledgeType,
		KnowledgeSource:     input.KnowledgeSource,
		Confidentiality:     input.Confidentiality,
		Abstract:            input.Abstract,
		Content:             input.Content,
		Detection:           input.Detection,
		Mitigations:         input.Mitigations,
		Recommendations:     input.Recommendations,
		Directory:           input.Directory,
		Techniques:          input.Techniques,
		Tactics:             input.Tactics,
		UsedExploits:        input.UsedExploits,
		Alias:               input.Alias,
		VulType:             input.VulType,
		Affiliation:         input.Affiliation,
		UsedTools:           input.UsedTools,
		StrategicCapability: input.StrategicCapability,
		FirstActivity:       input.FirstActivity,
		LatestActivity:      input.LatestActivity,
		TargetedGeography:   input.TargetedGeography,
		TimeLine:            input.TimeLine,
		Scenario:            input.Scenario,
		Motivations:         input.Motivations,
		TargetedIndustry:    input.TargetedIndustry,
		Preparation:         input.Preparation,
		Alert:               input.Alert,
		Analysis:            input.Analysis,
		Traces:              input.Traces,
		Containment:         input.Containment,
		Eradication:         input.Eradication,
		Recovery:            input.Recovery,
		FollowUp:            input.FollowUp,
		DisposalProcess:     input.DisposalProcess,
		Cases:               input.Cases,
		Cve:                 input.Cve,
		Cnnvd:               input.Cnnvd,
		Cwd:                 input.Cwd,
		Cvss:                input.Cvss,
		Bugtraq:             input.Bugtraq,
		CvssStr:             input.CvssStr,
		Msf:                 input.Msf,
		Exploitdb:           input.Exploitdb,
		IsExp:               input.IsExp,
		Vendor:              input.Vendor,
		AppType:             input.AppType,
		Consequence:         input.Consequence,
		FingerPrint:         input.FingerPrint,
		RevisionDate:        input.RevisionDate,
		Products:            input.Products,
		Reference:           input.Reference,
		Author:              input.Author,
		UID:                 input.UID,
		SubTechniquesID:     input.SubTechniquesID,
		Platforms:           input.Platforms,
		AffectedVerison:     input.AffectedVerison,
		ThreatSeverity:      input.ThreatSeverity,
		Solution:            input.Solution,
		Cnvd:                input.Cnvd,
		Cwe:                 input.Cwe,
		AppName:             input.AppName,
		OrganizationIds:     input.OrganizationIds,
		IoC:                 input.IoC,
		TiName:              input.TiName,
		InputParameters:     input.InputParameters,
		OutputParameters:    input.OutputParameters,
		Success:             true, // 设置默认值
		Message:             "Created successfully",
	}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &doc, nil
}

func (r *mutationResolver) UpdateKnowledge(ctx context.Context, id string, input model.NewKnowledge) (*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":               input.Title,
			"tags":                input.Tags,
			"techniquesId":        input.TechniquesID,
			"tacticsId":           input.TacticsID,
			"knowledgeType":       input.KnowledgeType,
			"knowledgeSource":     input.KnowledgeSource,
			"confidentiality":     input.Confidentiality,
			"abstract":            input.Abstract,
			"content":             input.Content,
			"detection":           input.Detection,
			"mitigations":         input.Mitigations,
			"recommendations":     input.Recommendations,
			"directory":           input.Directory,
			"techniques":          input.Techniques,
			"tactics":             input.Tactics,
			"usedExploits":        input.UsedExploits,
			"alias":               input.Alias,
			"vulType":             input.VulType,
			"affiliation":         input.Affiliation,
			"usedTools":           input.UsedTools,
			"strategicCapability": input.StrategicCapability,
			"firstActivity":       input.FirstActivity,
			"latestActivity":      input.LatestActivity,
			"targetedGeography":   input.TargetedGeography,
			"timeLine":            input.TimeLine,
			"scenario":            input.Scenario,
			"motivations":         input.Motivations,
			"targetedIndustry":    input.TargetedIndustry,
			"preparation":         input.Preparation,
			"alert":               input.Alert,
			"analysis":            input.Analysis,
			"traces":              input.Traces,
			"containment":         input.Containment,
			"eradication":         input.Eradication,
			"recovery":            input.Recovery,
			"followUp":            input.FollowUp,
			"disposalProcess":     input.DisposalProcess,
			"cases":               input.Cases,
			"cve":                 input.Cve,
			"cnnvd":               input.Cnnvd,
			"cwd":                 input.Cwd,
			"cvss":                input.Cvss,
			"bugtraq":             input.Bugtraq,
			"cvssStr":             input.CvssStr,
			"msf":                 input.Msf,
			"exploitdb":           input.Exploitdb,
			"isExp":               input.IsExp,
			"vendor":              input.Vendor,
			"appType":             input.AppType,
			"consequence":         input.Consequence,
			"fingerPrint":         input.FingerPrint,
			"revisionDate":        input.RevisionDate,
			"products":            input.Products,
			"reference":           input.Reference,
			"author":              input.Author,
			"subTechniquesId":     input.SubTechniquesID,
			"platforms":           input.Platforms,
			"affectedVersion":     input.AffectedVerison,
			"threatSeverity":      input.ThreatSeverity,
			"solution":            input.Solution,
			"cnvd":                input.Cnvd,
			"cwe":                 input.Cwe,
			"appName":             input.AppName,
			"organizationIds":     input.OrganizationIds,
			"IoC":                 input.IoC,
			"TiName":              input.TiName,
			"inputParameters":     input.InputParameters,
			"outputParameters":    input.OutputParameters,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var result model.Knowledge
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteKnowledge 实现
func (r *mutationResolver) DeleteKnowledge(ctx context.Context, id string) (*model.DeletionStatus, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"_id": id} //链接_id
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return &model.DeletionStatus{Success: false, Message: err.Error()}, err
	}

	if deleteResult.DeletedCount == 0 {
		return &model.DeletionStatus{Success: false, Message: "No document found with that ID"}, nil
	}

	return &model.DeletionStatus{Success: true, Message: "Deleted successfully"}, nil
}

// SearchByKnowledgeType 实现
func (r *queryResolver) SearchByKnowledgeType(ctx context.Context, typeArg []string, nums int) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"knowledgeType": bson.M{"$in": typeArg}}

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// MitreByTacticsID 实现
func (r *queryResolver) MitreByTacticsID(ctx context.Context, tacticsID []string) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	// Ensure tacticsID is always an array
	if len(tacticsID) == 0 {
		return nil, errors.New("tacticsID must be provided")
	}

	filter := bson.M{"tacticsId": bson.M{"$in": tacticsID}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// MitreByTechniquesID 实现
func (r *queryResolver) MitreByTechniquesID(ctx context.Context, techniquesID []string) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	if len(techniquesID) == 0 {
		return nil, errors.New("tacticsID must be provided")
	}

	filter := bson.M{"techniquesId": bson.M{"$in": techniquesID}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// MitreBySubTechniquesID 实现
func (r *queryResolver) MitreBySubTechniquesID(ctx context.Context, subTechniquesID []string) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"techniquesId": bson.M{"$in": subTechniquesID}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Search 实现
// func (r *queryResolver) Search(ctx context.Context, where *model.KnowledgeFilter, keyword []string, nums int) ([]*model.Knowledge, error) {
// 	var (
// 		filterList        []bson.M
// 		handlerFilter     func(*model.KnowledgeFilter) (bson.M, error)
// 		handleFilterError func(error, []string, int) ([]*model.Knowledge, error)
// 	)

// 	handlerFilter = func(where *model.KnowledgeFilter) (bson.M, error) {
// 		filter := bson.M{}
// 		whereMp := map[string]interface{}{}

// 		whereJson, err := bson.Marshal(where)
// 		if err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}

// 		if err := bson.Unmarshal(whereJson, &whereMp); err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}
// 		delete(whereMp, "AND")

// 		keyList := []string{}
// 		for key, value := range whereMp {
// 			if value != nil && key != "AND" {
// 				keyList = append(keyList, key)
// 			}
// 		}

// 		length := len(keyList)
// 		if length <= 0 {
// 			return nil, errors.New("0")
// 		}

// 		if length > 1 {
// 			return nil, errors.New("1")
// 		} else if len(keyList) < 1 {
// 			return nil, errors.New("2")
// 		}

// 		for i := 0; i < length; i++ {
// 			if util.IsListField(keyList[i]) {
// 				list, ok := whereMp[keyList[i]].(primitive.A)
// 				if !ok {
// 					return nil, fmt.Errorf("field %s is not a list", keyList[i])
// 				}
// 				interfaceList := []interface{}(list)
// 				filter[keyList[i]] = bson.M{"$in": interfaceList}
// 				continue
// 			}
// 			filter[keyList[i]] = bson.M{"$regex": whereMp[keyList[i]], "$options": "i"}
// 		}

// 		return filter, nil
// 	}

// 	handleFilterError = func(err error, keyword []string, nums int) ([]*model.Knowledge, error) {
// 		errorMessage := err.Error()
// 		switch errorMessage {
// 		case "0":
// 			return r.SearchByKeyword(ctx, []string{}, keyword, nums)
// 		case "1":
// 			return []*model.Knowledge{{Success: false, Message: "error: 输入字段过多，过滤器字段在AND中输入"}}, errors.New("输入字段过多，过滤器字段在AND中输入")
// 		case "2":
// 			return []*model.Knowledge{{Success: false, Message: "error: 未输入字段，请输入搜索一个字段及内容"}}, errors.New("未输入字段，请输入搜索一个字段及内容")
// 		default:
// 			log.Println(err)
// 			return nil, err
// 		}
// 	}

// 	dqlFilter, err := handlerFilter(where)
// 	if err != nil {
// 		return handleFilterError(err, keyword, nums)
// 	}

// 	filterList = append(filterList, dqlFilter)
// 	collection := database.GetCollection("knowledge")

// 	findOptions := options.Find()
// 	if nums > 0 {
// 		findOptions.SetLimit(int64(nums))
// 	}

// 	if len(where.AND) > 0 {
// 		for _, kf := range where.AND {
// 			andFilter, err := handlerFilter(kf)
// 			if err != nil {
// 				return handleFilterError(err, keyword, nums)
// 			}
// 			filterList = append(filterList, andFilter)
// 		}
// 	}

// 	finalFilter := bson.M{"$and": filterList}

// 	if len(keyword) > 0 {
// 		orConditions := []bson.M{}
// 		for _, key := range keyword {
// 			orConditions = append(orConditions, bson.M{"id": key})
// 			orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
// 			orConditions = append(orConditions, bson.M{"tags": key})
// 			orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
// 		}
// 		finalFilter["$or"] = orConditions
// 	}

// 	cursor, err := collection.Find(ctx, finalFilter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []*model.Knowledge
// 	for cursor.Next(ctx) {
// 		var result model.Knowledge
// 		err := cursor.Decode(&result)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, &result)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

//		return results, nil
//	}
//
// 修改后的第一个版本
// func (r *queryResolver) Search(ctx context.Context, where *model.KnowledgeFilter, keyword []string, nums int) ([]*model.Knowledge, error) {
// 	var filter bson.M

// 	// Build filter based on KnowledgeFilter fields
// 	whereJson, err := bson.Marshal(where)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	if err := bson.Unmarshal(whereJson, &filter); err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}

// 	// Add keyword search
// 	if len(keyword) > 0 {
// 		orConditions := []bson.M{}
// 		for _, key := range keyword {
// 			orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
// 			orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
// 			orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
// 		}
// 		filter["$or"] = orConditions
// 	}

// 	collection := database.GetCollection("knowledge")

// 	findOptions := options.Find()
// 	if nums > 0 {
// 		findOptions.SetLimit(int64(nums))
// 	}

// 	cursor, err := collection.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []*model.Knowledge
// 	for cursor.Next(ctx) {
// 		var result model.Knowledge
// 		err := cursor.Decode(&result)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, &result)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

//		return results, nil
//	}
//
// 修改二
func (r *queryResolver) Search(ctx context.Context, where *model.KnowledgeFilter, keyword []string, authors []string, nums int, nodedict string) ([]*model.Knowledge, error) {
	var filter bson.M

	// Build filter based on KnowledgeFilter fields
	whereJson, err := bson.Marshal(where)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := bson.Unmarshal(whereJson, &filter); err != nil {
		log.Println(err)
		return nil, err
	}

	// Add keyword search
	if len(keyword) > 0 {
		orConditions := []bson.M{}

		for _, key := range keyword {
			// orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
			// orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
			// orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
			switch nodedict {
			case "title":
				orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
			case "content":
				orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
			case "abstract":
				orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
			default:
				orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
				orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
				orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
			}
		}
		filter["$or"] = orConditions
	}

	// Add author search
	if len(authors) > 0 {
		filter["author"] = bson.M{"$all": authors}
	}

	if len(where.Tags) > 0 {
		filter["tags"] = bson.M{"$in": where.Tags}
	}

	collection := database.GetCollection("knowledge")

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// func (r *queryResolver) Search(ctx context.Context, where *model.KnowledgeFilter, keyword []string, authors []string, nums int, nodedict string) ([]*model.Knowledge, error) {
// 	var (
// 		filterList        []bson.M
// 		handlerFilter     func(*model.KnowledgeFilter) (bson.M, error)
// 		handleFilterError func(error, []string, int) ([]*model.Knowledge, error)
// 	)

// 	handlerFilter = func(where *model.KnowledgeFilter) (bson.M, error) {
// 		filter := bson.M{}
// 		whereMp := map[string]interface{}{}

// 		whereJson, err := bson.Marshal(where)
// 		if err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}

// 		if err := bson.Unmarshal(whereJson, &whereMp); err != nil {
// 			log.Println(err)
// 			return nil, err
// 		}
// 		delete(whereMp, "AND")

// 		keyList := []string{}
// 		for key, value := range whereMp {
// 			if value != nil && key != "AND" {
// 				keyList = append(keyList, key)
// 			}
// 		}

// 		length := len(keyList)
// 		if length <= 0 {
// 			return nil, errors.New("0")
// 		}

// 		if length > 1 {
// 			return nil, errors.New("1")
// 		} else if len(keyList) < 1 {
// 			return nil, errors.New("2")
// 		}

// 		for i := 0; i < length; i++ {
// 			if util.IsListField(keyList[i]) {
// 				list, ok := whereMp[keyList[i]].(primitive.A)
// 				if !ok {
// 					return nil, fmt.Errorf("field %s is not a list", keyList[i])
// 				}
// 				interfaceList := []interface{}(list)
// 				filter[keyList[i]] = bson.M{"$in": interfaceList}
// 				continue
// 			}
// 			filter[keyList[i]] = bson.M{"$regex": whereMp[keyList[i]], "$options": "i"}
// 		}

// 		return filter, nil
// 	}

// 	handleFilterError = func(err error, keyword []string, nums int) ([]*model.Knowledge, error) {
// 		errorMessage := err.Error()
// 		switch errorMessage {
// 		case "0":
// 			return r.SearchByKeyword(ctx, []string{}, keyword, nums)
// 		case "1":
// 			return []*model.Knowledge{{Success: false, Message: "error: 输入字段过多，过滤器字段在AND中输入"}}, errors.New("输入字段过多，过滤器字段在AND中输入")
// 		case "2":
// 			return []*model.Knowledge{{Success: false, Message: "error: 未输入字段，请输入搜索一个字段及内容"}}, errors.New("未输入字段，请输入搜索一个字段及内容")
// 		default:
// 			log.Println(err)
// 			return nil, err
// 		}
// 	}

// 	dqlFilter, err := handlerFilter(where)
// 	if err != nil {
// 		return handleFilterError(err, keyword, nums)
// 	}

// 	filterList = append(filterList, dqlFilter)
// 	collection := database.GetCollection("knowledge")

// 	findOptions := options.Find()
// 	if nums > 0 {
// 		findOptions.SetLimit(int64(nums))
// 	}

// 	if len(where.AND) > 0 {
// 		for _, kf := range where.AND {
// 			andFilter, err := handlerFilter(kf)
// 			if err != nil {
// 				return handleFilterError(err, keyword, nums)
// 			}
// 			filterList = append(filterList, andFilter)
// 		}
// 	}

// 	finalFilter := bson.M{"$and": filterList}

// 	if len(keyword) > 0 {
// 		orConditions := []bson.M{}
// 		for _, key := range keyword {
// 			switch nodedict {
// 			case "title":
// 				orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
// 			case "content":
// 				orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
// 			case "abstract":
// 				orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
// 			default:
// 				orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
// 				orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
// 				orConditions = append(orConditions, bson.M{"abstract": bson.M{"$regex": key, "$options": "i"}})
// 			}
// 		}
// 		finalFilter["$or"] = orConditions
// 	}

// 	if len(authors) > 0 {
// 		finalFilter["author"] = bson.M{"$all": authors}
// 	}

// 	cursor, err := collection.Find(ctx, finalFilter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var results []*model.Knowledge
// 	for cursor.Next(ctx) {
// 		var result model.Knowledge
// 		err := cursor.Decode(&result)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, &result)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	return results, nil
// }

// SearchByKeyword 实现
func (r *queryResolver) SearchByKeyword(ctx context.Context, typeArg []string, keyword []string, nums int) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{}
	if len(typeArg) > 0 {
		filter["knowledgeType"] = bson.M{"$in": typeArg}
	}
	if len(keyword) > 0 {
		orConditions := []bson.M{}
		for _, key := range keyword { // 搜索关键字的字段
			orConditions = append(orConditions, bson.M{"id": key})
			orConditions = append(orConditions, bson.M{"title": bson.M{"$regex": key, "$options": "i"}})
			orConditions = append(orConditions, bson.M{"tags": key})
			orConditions = append(orConditions, bson.M{"content": bson.M{"$regex": key, "$options": "i"}})
		}
		filter["$or"] = orConditions
	}

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// SearchByTagsWithType 实现
func (r *queryResolver) SearchByTagsWithType(ctx context.Context, typeArg []string, tags []string, nums int) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{}
	if len(typeArg) > 0 {
		filter["knowledgeType"] = bson.M{"$in": typeArg}
	}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$all": tags}
	}

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// SearchByContent 实现
func (r *queryResolver) SearchByContent(ctx context.Context, typeArg []string, keyword string, nums int) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{}
	if len(typeArg) > 0 {
		filter["knowledgeType"] = bson.M{"$in": typeArg}
	}
	if keyword != "" {
		filter["content"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// SearchByTitle 实现
func (r *queryResolver) SearchByTitle(ctx context.Context, title string, nums int) ([]*model.Knowledge, error) {
	collection := database.GetCollection("knowledge")

	filter := bson.M{"title": bson.M{"$regex": title, "$options": "i"}}

	findOptions := options.Find()
	if nums > 0 {
		findOptions.SetLimit(int64(nums))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Knowledge
	for cursor.Next(ctx) {
		var result model.Knowledge
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
