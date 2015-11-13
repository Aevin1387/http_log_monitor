package main

import (
  "fmt"
  "regexp"
  "strconv"
  "strings"
  "time"
)

type LineData struct {
  Date    time.Time

  RequestStr    string
  RequestMethod string
  SectionStr    string
  RequestProtocol string

  Status     string
  ContentLen int

  RemoteHost string
  Rfc931     string
  AuthUser   string
}

type RegexMap struct {
  regex string
  field string
}

const (
  W3CDate = "02/Jan/2006:15:04:05 -0700"
)

func startParser(inputChan chan string, outputChan chan LineData, logFormat string) {
  formatsMaps := parseFormatToRegexMaps(logFormat)
  formatsRegex := parseRegexMapsToRegex(formatsMaps)

  for {
    select {
      case line := <- inputChan:
        outputChan <- parseLine(line, formatsMaps, formatsRegex)
    }
  }
}

func parseLine(line string, formatsMaps []RegexMap, formatsRegex *regexp.Regexp) LineData {
  match := formatsRegex.FindStringSubmatch(line)
  data := LineData{}

  for index, value := range match {
    if(index == 0) {
      continue
    }

    mapping := formatsMaps[index - 1]
    switch mapping.field {
    case "host":
      data.RemoteHost = value
    case "ident":
    case "user":
      data.AuthUser = value
    case "date":
      date, err := time.Parse(W3CDate,value)
      if(err != nil) {
        fmt.Println(err)
      }
      data.Date = date
    case "request":
      requests := strings.Split(value, " ")
      data.RequestStr = value
      data.RequestMethod = requests[0]

      location := requests[1]
      locations := strings.Split(location, "/")
      if(len(locations) > 1) { // First / considered blank
        data.SectionStr = "/" + locations[1]
      } else {
        data.SectionStr = "/"
      }

      data.RequestProtocol = requests[2]
    case "status":
      data.Status = value
    case "size":
      length, _ := strconv.Atoi(value)
      data.ContentLen = length
    default:
    }
  }

  return data
}


func parseFormatToRegexMaps(logFormat string) []RegexMap {
  formats := strings.Split(logFormat, " ")
  formatRegexes := make([]RegexMap, 0)

  // "%h %l %u %t \"%r\" %>s %b"
  for _, format := range formats {
    switch format {
    case "%h":
      formatRegexes = append(formatRegexes, RegexMap{ field: "host", regex: "([\\S]+)" })
    case "%l":
      formatRegexes = append(formatRegexes, RegexMap{ field: "ident", regex: "([\\-\\w]+)" })
    case "%u":
      formatRegexes = append(formatRegexes, RegexMap{ field: "user", regex: "([\\-\\w]+)" })
    case "%t":
      formatRegexes = append(formatRegexes, RegexMap{ field: "date", regex: "\\[(\\d{2}\\/\\w{3}\\/\\d{4}:\\d{2}:\\d{2}:\\d{2} ?-?\\d{4}?)\\]" })
    case "\"%r\"":
      formatRegexes = append(formatRegexes, RegexMap{ field: "request", regex: "\"(.*)\"" })
    case "%>s":
      formatRegexes = append(formatRegexes, RegexMap{ field: "status", regex: "(\\d*-?)" })
    case "%b":
      fallthrough
    case "%B":
      formatRegexes = append(formatRegexes, RegexMap{ field: "size", regex: "(\\d*-?)" })
    }
  }

  return formatRegexes
}

func parseRegexMapsToRegex(formatsMaps []RegexMap) *regexp.Regexp {
  regexs := make([]string, 0)

  for _, regexMap := range formatsMaps {
    regexs = append(regexs, regexMap.regex)
  }
  regexStr := strings.Join(regexs, " ")

  regexStr = "^" + regexStr
  r := regexp.MustCompile(regexStr)

  return r
}
