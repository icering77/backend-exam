package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Employee struct {
	name string
}

type Item1 struct {
	id       int
	itemType int
}

type Item2 struct {
	id       int
	itemType int
}

type Item3 struct {
	id       int
	itemType int
}

type Item interface {
	// Process 這是一個耗時操作
	getType() int
	getId() int
	Process() int
}

func (item *Item1) Process() int {
	time.Sleep(time.Second * 1)
	return 1
}

func (item *Item2) Process() int {
	time.Sleep(time.Second * 2)
	return 2
}

func (item *Item3) Process() int {
	time.Sleep(time.Second * 3)
	return 3
}

func (item *Item1) getType() int {
	return item.itemType
}

func (item *Item1) getId() int {
	return item.id
}

func (item *Item2) getType() int {
	return item.itemType
}

func (item *Item2) getId() int {
	return item.id
}

func (item *Item3) getType() int {
	return item.itemType
}

func (item *Item3) getId() int {
	return item.id
}

func processItem(emp Employee, items <-chan Item, done chan<- bool) {
	for item := range items {
		fmt.Println(emp.name, "開始處理商品類型:", item.getType(), "商品編號:", item.getId())
		fmt.Println(emp.name, "完成處理商品類型:", item.getType(), "商品編號:", item.getId(), "耗時:", item.Process(), "秒")
	}
	done <- true
}

func main() {
	// 建立商品並打亂
	itemsArr := make([]Item, 0, 30)
	for i := 0; i < 10; i++ {
		itemsArr = append(itemsArr, &Item1{
			id:       i,
			itemType: 1,
		})
	}
	for i := 0; i < 10; i++ {
		itemsArr = append(itemsArr, &Item2{
			id:       i,
			itemType: 2,
		})
	}
	for i := 0; i < 10; i++ {
		itemsArr = append(itemsArr, &Item3{
			id:       i,
			itemType: 3,
		})
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(itemsArr), func(i, j int) {
		itemsArr[i], itemsArr[j] = itemsArr[j], itemsArr[i]
	})

	// 建立員工
	Employees := make([]Employee, 0, 5)
	for i := 0; i < 5; i++ {
		Employees = append(Employees, Employee{
			name: fmt.Sprintf("員工%d", i+1),
		})
	}

	itemChan := make(chan Item, len(itemsArr))
	done := make(chan bool, len(Employees))

	// 啟動五個員工 goroutine
	for _, emp := range Employees {
		go processItem(emp, itemChan, done)
	}

	// 將商品放入 channel
	for _, item := range itemsArr {
		itemChan <- item
	}
	close(itemChan) // 發完所有商品

	// 等待所有員工完成
	for i := 0; i < len(Employees); i++ {
		<-done
	}

	fmt.Println("所有商品處理完成")

}
