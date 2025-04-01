package test

import (
	"chg/internal/api/imagesearch"
	"chg/internal/api/imagesearch/fetcher"
	"fmt"
	"testing"
)

func TestGetImagePageURL(t *testing.T) {
	// 测试用例：传入一个图片URL
	imageURL := "https://cloudhivegallery-1348386678.cos.ap-guangzhou.myqcloud.com/public/557094266849998849/2025-03-28_3d980ca87c556e15.webp"

	// 调用函数
	result, err := fetcher.GetImagePageURL(imageURL)

	// 断言结果是否符合预期
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}
	//输出结果
	t.Logf("Result: %s", result)
	fmt.Printf("Result: %s\n", result)
}

func TestGetImageFirstURL(t *testing.T) {
	// 测试用例：传入一个图片URL
	imageURL := "https://graph.baidu.com/s?card_key=&entrance=GENERAL&extUiData%5BisLogoShow%5D=1&f=all&isLogoShow=1&session_id=4916634739113792653&sign=126e8ec9ae221b697484601743440911&tpl_from=pc"

	// 调用函数
	result, err := fetcher.GetImageFirstURL(imageURL)

	// 断言结果是否符合预期
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == "" {
		t.Errorf("Expected non-empty result, got empty string")
	}
	//输出结果
	t.Logf("Result: %s", result)
}

func TestGetImageList(t *testing.T) {
	// 测试用例：传入一个图片URL
	imageURL := "https://graph.baidu.com/ajax/pcsimi?carousel=503&entrance=GENERAL&extUiData%5BisLogoShow%5D=1&inspire=general_pc&limit=30&next=2&render_type=card&session_id=4916634739113792653&sign=126e8ec9ae221b697484601743440911&tk=e3d40&tpl_from=pc"

	// 调用函数
	result, err := fetcher.GetImageList(imageURL)
	// 断言结果是否符合预期
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Errorf("Expected non-empty result, got empty string")
	}
	//输出结果
	t.Logf("Result: %s", result)
}

func TestSearchImage(t *testing.T) {
	// 测试用例：传入一个图片URL
	imageURL := "https://cloudhivegallery-1348386678.cos.ap-guangzhou.myqcloud.com/public/557094266849998849/2025-03-28_3d980ca87c556e15.webp"

	// 调用函数
	result, err := imagesearch.SearchImage(imageURL)

	// 断言结果是否符合预期
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Errorf("Expected non-empty result, got empty string")
	}
	//输出结果
	t.Logf("Result: %s", result)
}
