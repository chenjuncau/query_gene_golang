// read file into the map
// https://golang.cafe/blog/golang-zip-file-example.html
// https://earthly.dev/blog/golang-zip-files/
// https://askgolang.com/golang-zip-file/
// This is working.
package main

import (
	"bufio"
	// "encoding/csv"
	"fmt"
	// "io"
	"log"
	"os"
	"strings"
	"strconv"
)

func containstring(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
// This is the first part, Open the file and put them into map hash table 
// This pare can be change read the row #CHROM. and get the position as the index. and then query them out.easy here.
	fileFormat := "VCF"

	// read the chr:postion file and then subset it. 
    readFile, err := os.Open("/data/Jun/galGal7/genes/gene.ensembl.list2")    // genes.txt  gene.ensembl.list (mouse) chicken.gene.ensembl.list (chicken)
    var mRep []string
    if err != nil {
        fmt.Println(err)
    }
    fileScanner := bufio.NewScanner(readFile)

//    fileScanner.Split(bufio.ScanLines)

    for fileScanner.Scan() {
      
//        fmt.Println(fileScanner.Text())
        record :=strings.Split(fileScanner.Text(), "\t")
	recordMerge :=strings.Join(record, ":")
	mRep = append(mRep, recordMerge)
}

    readFile.Close()


	// // Iterate through the records
	// for {
		// // Read each record from csv
		// record, err := r.Read()
		// if err == io.EOF {
			// break
		// }
		// if err != nil {
			// log.Fatal(err)
		// }
		// fmt.Printf("Question: %s Answer %s\n", record[0], record[1])
		// mRep[record[0]] = record[1]
// //		fmt.Println("This is the map check")
// //		fmt.Println("map:", mRep)
// //		fmt.Println("map value:", mRep[record[0]])
	// }
	
// This is second part, processing the VCF file, AB format file and sepplemental file.	

    // f, err := os.Create("/work/51/Jun/RP_name/Test2output_geno_summary.vcf")
    f, err := os.Create("/data/Jun/galGal7/genes/gg7_deleterious_variants.reformated.vcf")  //output file
    check(err)
    defer f.Close()
	
	// file, err := os.Open("/work/51/Jun/RP_name/Test2.vcf")  // input vcf file.
	file, err := os.Open("/data/Jun/galGal7/gg7_gs_data.recode.vcf")  // input vcf file.
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 20480*10240)
	lineNumber := 0
	// temp1 :=""
	printNum :=0
	snpCount :=0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text() // Read the current line
		homRef := 0
		het := 0
		homAlt := 0
		total :=0
		// Process the line here
		if fileFormat=="VCF" {
			if strings.HasPrefix(line, "#CHROM") {  //vcf file
			// if lineNumber==10 {  // AB format file
	//			fmt.Println("Processing line:", lineNumber)
				printNum = lineNumber
				strings.Replace(line, "#CHROM", "CHROM",-1)
				_, err := f.WriteString(line+"\t"+"geneName" + "\t" + "snpIndex" + "\t" + "hom_ref"+"\t"+"het"+"\t"+"hom_alt"+"\t"+"total"+"\n") // add other condition, so you can change it.  
				check(err)				
			}
			if (lineNumber>printNum) && (printNum !=0) {  //vcf file
				snpCount++
	//		    fmt.Println("Line Number is :", lineNumber)
				sampledata :=strings.Split(line, "\t")
// adding the filter chr and pos here. to filter.
// tempJoin :=sampledata[0] + ":" + sampledata[1]
// filter by genes.
// tempJoin :=sampledata[0] + ":" + sampledata[1]
//if (!containstring(mRep, tempJoin))  {
// continue 
//}

found := false
geneName := "unknown"
  for _, str := range mRep {
    if strings.Contains(sampledata[7], str) {
          found = true
	  fmt.Println("Processing gene name:", str)
	  geneName=str
          break
     }
  }

if !(found)  {
 continue
 }      

// end of filter. 
				for i, data := range sampledata {
					if i >= 9 {
						temp1 :=strings.Split(data, ":") 
					   // sampledata[i]=temp1[0]   // only genotype
//					      sampledata[i]=strings.Join(temp1[0:2], ":") // genotype and readdepth // this is rapid vcf file.
                                             sampledata[i]=temp1[0] + ":" + temp1[2]  // interval-bio vcf genotype and readdepth
					    // sampledata[i]=temp1[1]  // only total readdepth
					    // sampledata[i]=temp1[2]  // only seperate readdepth
						// if strings.contain(temp1[0],"0/0") {homRef++}  // second method
					}
					// if _, ok := mRep[data]; ok {			
						// sampledata[i]=mRep[data]
					// }
				}
	// refref  altalt  het     total   refref_frq      alt_frq het_frq singleSNP			
				homRef=strings.Count(line,"\t0/0:")+strings.Count(line,"\t0|0:")
				het=strings.Count(line,"\t0/1:")+strings.Count(line,"\t0|1:")+strings.Count(line,"\t1/0:")+strings.Count(line,"\t1|0:")
				homAlt=strings.Count(line,"\t1/1:")+strings.Count(line,"\t1|1:")
				total=homRef+het+homAlt
				sampledata = append(sampledata,geneName,strconv.Itoa(snpCount),strconv.Itoa(homRef),strconv.Itoa(het),strconv.Itoa(homAlt),strconv.Itoa(total))
				data2 :=strings.Join(sampledata, "\t")
				// _, err = f.WriteString(data2 +"\t"+homRef+"\t"+strings(het)+"\t"+strings(homAlt)+"\t"+strings(total)+"\n")
				_, err = f.WriteString(data2 +"\n")

			} 
			// else {
				// _, err := f.WriteString(line + "\n") // add other condition, so you can change it.  
				// check(err)
			// }
		}
		
	}
    f.Sync()
	if err := scanner.Err(); err != nil {
    log.Fatalf("something bad happened in the line %v: %v", lineNumber, err)
	}	
	
}


// if(the first value is CHROM) then 
// https://devmarkpro.com/working-big-files-golang
// output := strings.Replace(input, "|", "\n", -1)
// https://gobyexample.com/string-functions			
		// fmt.Println("Processing line:", line)
		// fmt.Printf("wrote %d bytes\n", n3)
   
	// if err := scanner.Err(); err != nil {
		// fmt.Println("Error reading file:", err)
	// }	
