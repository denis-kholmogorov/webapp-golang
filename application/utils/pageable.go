package utils

type Pageable interface {
	GetPage() int
}

func GetPageNumber(pageable Pageable) int {
	if pageable.GetPage() > 0 {
		return pageable.GetPage()
	} else {
		return 0
	}
}
