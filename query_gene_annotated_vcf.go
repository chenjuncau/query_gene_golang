# Jun Chen
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

func containString(elems []string, v string) bool {
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
	fileFormat := "VCF"

	readFile, err := os.Open("gene.ensembl.list")
    var mRep []string
    if err != nil {
        fmt.Println(err)
    }
    fileScanner := bufio.NewScanner(readFile)

    for fileScanner.Scan() {
        record := strings.Split(fileScanner.Text(), "\t")
	    recordMerge := strings.Join(record, ":")
	    mRep = append(mRep, recordMerge)
    }
    readFile.Close()

    f, err := os.Create("reformated.vcf")
    check(err)
    defer f.Close()

	file, err := os.Open("output.vcf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	buf := []byte{}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 20480*10240)
	lineNumber := 0
	printNum := 0
	snpCount := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		homRef := 0
		het := 0
		homAlt := 0
		total := 0

		if fileFormat == "VCF" {
			if strings.HasPrefix(line, "#CHROM") {
				printNum = lineNumber
				strings.Replace(line, "#CHROM", "CHROM", -1)
				_, err := f.WriteString(line + "\t" + "geneName" + "\t" + "snpIndex" + "\t" + "hom_ref" + "\t" + "het" + "\t" + "hom_alt" + "\t" + "total" + "\n")
				check(err)
			}
			if lineNumber > printNum && printNum != 0 {
				snpCount++
				sampledata := strings.Split(line, "\t")
				found := false
				geneName := "unknown"
  				for _, str := range mRep {
    				if strings.Contains(sampledata[7], str) {
        				found = true
	  				    fmt.Println("Processing gene name:", str)
	  				    geneName = str
        				break
    				}
  				}
  				if !found {
    				continue
  				}
				for i, data := range sampledata {
					if i >= 9 {
						temp1 := strings.Split(data, ":")
					    sampledata[i] = temp1[0] + ":" + temp1[2]
					}
				}
				homRef = strings.Count(line, "\t0/0:") + strings.Count(line, "\t0|0:")
				het = strings.Count(line, "\t0/1:") + strings.Count(line, "\t0|1:") + strings.Count(line, "\t1/0:") + strings.Count(line, "\t1|0:")
				homAlt = strings.Count(line, "\t1/1:") + strings.Count(line, "\t1|1:")
				total = homRef + het + homAlt
				sampledata = append(sampledata, geneName, strconv.Itoa(snpCount), strconv.Itoa(homRef), strconv.Itoa(het), strconv.Itoa(homAlt), strconv.Itoa(total))
				data2 := strings.Join(sampledata, "\t")
				_, err = f.WriteString(data2 + "\n")
			}
		}
	}
    f.Sync()
	if err := scanner.Err(); err != nil {
        log.Fatalf("something bad happened in the line %v: %v", lineNumber, err)
	}
}
