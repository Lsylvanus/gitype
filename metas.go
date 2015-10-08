// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"

	"github.com/caixw/typing/core"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
	"github.com/issue9/orm/fetch"
)

// @api put /admin/api/tags/{id}/merge 将指定的标签合并到当前标签
// @apiGroup admin
//
// @apiRequest json
// @apiParam tags array 所有需要合并的标签ID列表。
// @apiExample json
// {"tags": [1,2,3] }
func adminPutTagMerge(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// @api get /admin/api/cats 获取所有分类信息
// @apiGroup admin
//
// @apiRequest json
// @apiheader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam cats array 所有分类的列表
func adminGetCats(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT m.{name},m.{title},m.{description},m.{id},m.{parent},m.{order},COUNT(r.{metaID}) AS {count}
			FROM #metas AS m
			LEFT JOIN #relationships AS r ON m.{id}=r.{metaID}
			WHERE m.{type}=?
			GROUP BY m.{id}
			ORDER BY m.{order} ASC`
	rows, err := db.Query(true, sql, models.MetaTypeCat)
	if err != nil {
		logs.Error("frontGetCats:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		logs.Error("frontGetCats:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"cats": maps}, nil)
}

// @api get /admin/api/tags 获取所有标签信息
// @apiGroup admin
//
// @apiRequest json
// @apiheader Authorization xxx
//
// @apiSuccess 200 OK
// @apiParam tags array 所有分类的列表
func adminGetTags(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT m.{name},m.{title},m.{description},m.{id},count(r.{metaID}) AS {count}
			FROM #metas AS m
			LEFT JOIN #relationships AS r ON m.{id}=r.{metaID}
			WHERE m.{type}=?
			GROUP BY m.{id}`
	rows, err := db.Query(true, sql, models.MetaTypeTag)
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	maps, err := fetch.MapString(false, rows)
	rows.Close()
	if err != nil {
		logs.Error("getTags:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusOK, map[string]interface{}{"tags": maps}, nil)
}

// @api patch /admin/api/{id}/order 修改某一分类的显示顺序
// @apiParam id int 分类的id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam order int 排序值
//
// @apiSuccess 204 No Content
func adminPatchCatOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}
	cat := &models.Meta{ID: id}
	if err := db.Select(cat); err != nil {
		logs.Error("adminPatchCatOrder:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	if cat.Type != models.MetaTypeCat {
		core.RenderJSON(w, http.StatusNotFound, nil, nil)
		return
	}

	o := &struct {
		Order int `json:"order"`
	}{}
	if !core.ReadJSON(w, r, o) {
		return
	}
	cat = &models.Meta{
		ID:    id,
		Order: o.Order,
	}
	if _, err := db.Update(cat); err != nil {
		logs.Error("adminPatchCatOrder:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// @api put /admin/api/tags/{id} 修改某id的标签内容
// @apiParam id int 需要修改的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 204 no content
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPutTag(w http.ResponseWriter, r *http.Request) {
	putMeta(w, r, models.MetaTypeTag)
}

// @api put /admin/api/cats/{id} 修改某id的分类内容
// @apiParam id int 需要修改的分类id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam parent int 父类
// @apiParam order int 排列顺序
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "parent": 5,
//     "order": 10,
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 200 ok
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPutCat(w http.ResponseWriter, r *http.Request) {
	putMeta(w, r, models.MetaTypeCat)
}

// @api post /admin/api/tags 添加新标签
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 201 created
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPostTag(w http.ResponseWriter, r *http.Request) {
	postMeta(w, r, models.MetaTypeTag)
}

// @api post /admin/api/cats 添加新的分类
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
// @apiParam name string 唯一名称
// @apiParam title string 显示的标题
// @apiParam parent int 父类
// @apiParam order int 排列顺序
// @apiParam description string 描述信息，可以是html
// @apiExample json
// {
//     "name": "tag-1",
//     "title":"标签1",
//     "parent": 5,
//     "order": 10,
//     "description": "<h1>desc</h1>"
// }
//
// @apiSuccess 201 created
// @apiError 400 bad request
// @apiParam message string 错误信息
// @apiParam detail array 说细的错误信息，用于描述哪个字段有错
// @apiExample json
// {
//     "message": "格式错误",
//     "detail":[
//         {"title":"不能包含特殊字符"},
//         {"name": "已经存在同名"}
//     ]
// }
func adminPostCat(w http.ResponseWriter, r *http.Request) {
	postMeta(w, r, models.MetaTypeCat)
}

// @api delete /admin/api/tags/{id} 删除该id的标签
// @apiParam id int 需要删除的标签id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 no content
func adminDeleteTag(w http.ResponseWriter, r *http.Request) {
	deleteMeta(w, r)
}

// @api put /admin/api/cats/{id} 删除该id的分类
// @apiParam id int 需要删除的分类id
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization xxx
//
// @apiSuccess 204 no content
func adminDeleteCat(w http.ResponseWriter, r *http.Request) {
	deleteMeta(w, r)
}

// 是否存在相同name的meta
func metaNameIsExists(m *models.Meta, typ int) (bool, error) {
	sql := db.Where("{name}=?", m.Name).And("{type}=?", typ).Table("#metas")
	maps, err := sql.SelectMapString(true, "id")
	if err != nil {
		return false, err
	}

	if len(maps) == 0 {
		return false, nil
	}
	if len(maps) > 1 {
		return true, nil
	}

	id, err := strconv.ParseInt(maps[0]["id"], 10, 64)
	println(maps[0]["id"])
	if err != nil {
		return false, err
	}
	return id != m.ID, nil
}

// @api get /admin/api/cats/{id} 获取指定id的分类内容
// @apiParam id int 分类的id
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam id int 分类的id
// @apiParam name string 分类的唯一名称，可能为空
// @apiParam title string 分类名称
// @apiParam order int 分类的排序值
// @apiParam description string 对分类的详细描述
func adminGetCat(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	m := &models.Meta{ID: id}
	if err := db.Select(m); err != nil {
		logs.Error("adminGetTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := &struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Order       int    `json:"order"`
		Description string `json:"description"`
	}{
		ID:          m.ID,
		Name:        m.Name,
		Title:       m.Title,
		Order:       m.Order,
		Description: m.Description,
	}
	core.RenderJSON(w, http.StatusOK, data, nil)
}

// @api get /admin/api/tags/{id} 获取指定id的标签内容
// @apiParam id int 标签的id
// @apiGroup admin
//
// @apiSuccess 200 OK
// @apiParam id int 标签的id
// @apiParam name string 标签的唯一名称，可能为空
// @apiParam title string 标签名称
// @apiParam description string 对标签的详细描述
func adminGetTag(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	m := &models.Meta{ID: id}
	if err := db.Select(m); err != nil {
		logs.Error("adminGetTag:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	data := &struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}{
		ID:          m.ID,
		Name:        m.Name,
		Title:       m.Title,
		Description: m.Description,
	}
	core.RenderJSON(w, http.StatusOK, data, nil)
}

// 供putCat和putTag调用
func putMeta(w http.ResponseWriter, r *http.Request, typ int) {
	m := &models.Meta{}
	if !core.ReadJSON(w, r, m) {
		return
	}

	var ok bool
	m.ID, ok = core.ParamID(w, r, "id")
	if !ok {
		return
	}

	exists, err := metaNameIsExists(m, typ)
	if err != nil {
		logs.Error("putMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	errs := &core.ErrorResult{Message: "格式错误"}
	if exists {
		errs.Detail["name"] = "已有同名字体段"
	}

	if len(m.Title) == 0 {
		errs.Detail["title"] = "标题不能为空"
	}

	if len(errs.Detail) > 0 {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}

	if _, err := db.Update(m); err != nil {
		logs.Error("putMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 供postCat和postTag调用
func postMeta(w http.ResponseWriter, r *http.Request, typ int) {
	m := &models.Meta{}
	if !core.ReadJSON(w, r, m) {
		return
	}

	errs := &core.ErrorResult{Message: "格式错误"}
	if m.ID > 0 {
		errs.Detail["id"] = "不允许的字段"
	}
	m.ID = 0

	exists, err := metaNameIsExists(m, typ)
	if err != nil {
		logs.Error("postMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if exists {
		errs.Detail["name"] = "已有同名字体段"
	}

	if len(m.Title) == 0 {
		errs.Detail["title"] = "标题不能为空"
	}

	if len(errs.Detail) > 0 {
		core.RenderJSON(w, http.StatusBadRequest, errs, nil)
		return
	}
	m.Type = typ

	if _, err := db.Insert(m); err != nil {
		logs.Error("postMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
	core.RenderJSON(w, http.StatusCreated, "{}", nil)
}

// 删除meta数据，供deleteCat和deleteTag调用
func deleteMeta(w http.ResponseWriter, r *http.Request) {
	id, ok := core.ParamID(w, r, "id")
	if !ok {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logs.Error("deleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if _, err := tx.Delete(&models.Meta{ID: id}); err != nil {
		logs.Error("deleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	// 删除与之对应的关联数据。
	sql := "DELETE FROM #relationships WHERE {MetaID}=?"
	if _, err := tx.Exec(true, sql, id); err != nil {
		logs.Error("deleteMeta:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	if err := tx.Commit(); err != nil {
		logs.Error("deleteMeta:", err)
		tx.Rollback()
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}

	core.RenderJSON(w, http.StatusNoContent, nil, nil)
}

// 获取与某post相关联的标签或是分类
func getPostMetas(postID int64, mtype int) ([]int64, error) {
	sql := `SELECT rs.{metaID} FROM #relationships AS rs
	LEFT JOIN #metas AS m ON m.{id}=rs.{metaID}
	WHERE rs.{postID}=? AND m.{type}=?`
	rows, err := db.Query(true, sql, postID, mtype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maps, err := fetch.ColumnString(false, "metaID", rows)
	if err != nil {
		return nil, err
	}

	ret := make([]int64, 0, len(maps))
	for _, v := range maps {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, num)
	}
	return ret, nil
}
