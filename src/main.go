package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type order struct {
	ID   uint
	name string
	deps []uint
}

var orders []order

func main() {
	getOrders()
	getDeps()
	writeOut()
}

func getOrders() {
	// open input file
	fi, err := os.Open("../data/orders.txt")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	// make a read buffer
	r := bufio.NewReader(fi)
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk into buf
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}

	orderSlice := strings.Split(string(buf), "\n")
	for o := range orderSlice {
		orderKV := strings.Split(orderSlice[o], ",")
		if _, err := strconv.Atoi(orderKV[0]); err == nil {
			orderID, _ := strconv.Atoi(orderKV[0])
			newOrder := order{ID: uint(orderID), name: orderKV[1]}
			orders = append(orders, newOrder)
		}
	}
}

func getDeps() {
	fi, err := os.Open("../data/dependencies.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	r := bufio.NewReader(fi)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}

	depSlice := strings.Split(string(buf), "\n")
	for d := range depSlice {
		depKV := strings.Split(depSlice[d], ",")
		if _, err := strconv.Atoi(depKV[0]); err == nil {
			orderID, _ := strconv.Atoi(depKV[0])
			orderDep, _ := strconv.Atoi(depKV[1])
			for o := range orders {
				if uint(orderID) == orders[o].ID {
					orders[o].deps = append(orders[o].deps, uint(orderDep))
				}
			}
		}
	}
}

func writeOut() {
	// open output file
	fo, err := os.Create("../output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// get list of orders that are not dependencies
	rootOrders := []order{}
	depdOnOrders := []order{}
	for _, o := range orders {
		for _, ao := range orders {
			for _, d := range ao.deps {
				if o.ID == d {
					depdOnOrders = append(depdOnOrders, o)
				}
			}
		}
		if orderInSlice(o, depdOnOrders) == false {
			rootOrders = append(rootOrders, o)
		}
	}

	// get rootOrders & deps in string format
	rootStrs := []string{}
	for _, o := range rootOrders {
		orderIDStr := strconv.Itoa(int(o.ID))
		orderStr := "Id: " + orderIDStr + ", Name: " + o.name
		rootStrs = append(rootStrs, orderStr)
	}
	fmt.Println("rootStrs: ", rootStrs)

	// print strings to file
	for _, p := range rootStrs {
		fmt.Fprintln(fo, p)
	}
}

func orderInSlice(o order, depdOrders []order) bool {
	for _, do := range depdOrders {
		if do.ID == o.ID {
			return true
		}
	}
	return false
}
