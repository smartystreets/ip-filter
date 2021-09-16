package IPFilter

/**
This file is to prepare the file of restricted IP addresses. We are using a radix sort to sort the data into groups
with similar node patterns.
*/

//Get file from somewhere containing all data types <- Need to verify that it will be converted to a slice of strings or
//something similar. Could potentially convert to digit, to sort easier.
//Filter through file to create radix tree

type RadixTree struct { //FIXME - need to add to this.
	length uint64
	max    uint64 //FIXME could be int?
}

//Calculate the longest string in the array to sort <- This will be the longest IP address
func (this *RadixTree) getMax(arr []int, n int) int {
	max := arr[0]
	for i := 0; i < n; i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}

func (this *RadixTree) countSort(arr []int, n, exp int) {
	var output []int  //Final output array
	var count [10]int //Could this be uint64?
	for i := 0; i < 10; i++{
		count[i] = 0
	}

	//Store count of occurrences in count[]
	for i := 0; i < n; i++ {
		count[(arr[i]/exp)%10]++
	}

	//Change count[i] so that count[i now contains the actual position of the digit in output
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	//Build the output array
	for i := n - 1; i >= 0; i-- {
		output[count[(arr[i]/exp)%10]-1] = arr[i]
		count[(arr[i]/exp)%10]--
	}

	// Copy the output array to the current array
	for i := 0; i < n; i++ {
		arr[i] = output[i]
	}
}

func (this *RadixTree) radixSort(arr []int, n int){
	m := this.getMax(arr, n)
	for exp := 1; m / exp > 0; exp *= 10{
		this.countSort(arr, n, exp)
	}
}
