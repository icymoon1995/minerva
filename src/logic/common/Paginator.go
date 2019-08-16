package logic

import "strconv"

/**
计算数据偏移量
@param page int 当前页码
@param perPage int 每页数据量
*/
func Offset(page int, perPage int) int {
	var offset int = (page - 1) * perPage
	if offset <= 0 {
		offset = 0
	}
	return offset
}

/**
分页处理
@param total int64 数据总数
@param currentPage int 当前页码
@param perPage int 每页的数量
@param list []interface{} 具体的数据
*/
func Paginate(total int64, currentPage int, perPage int, list []interface{}) map[string]interface{} {
	data := make(map[string]interface{})

	// 将int64 转为int
	intTotal, _ := strconv.Atoi(strconv.FormatInt(total, 10))

	// 计算分页相关数量
	var totalPage = intTotal / perPage
	if intTotal%perPage != 0 {
		totalPage += 1
	}

	data["total"] = intTotal
	data["current_page"] = currentPage
	data["per_page"] = perPage
	data["total_page"] = totalPage
	data["data"] = list
	return data
}
