package main

import (
  "time"
)

// SectionStat is a struct containing the
// statistics for a section starting at a period of time
type SectionStat struct {
  hitCount       int              // number of hits
  methodsCounts  map[string] int  // number of hits per method
  contentLengths []int            // size of each hit
  statusCounts   map[string] int  // number of hits per status
  hostCounts     map[string] int  // number of hits from each host
  startTime      time.Time        // time when statistics started
}

// SectionData is a struct containing
// the current and history of statistics for a
// section
type SectionData struct {
  name           string           // section name, ex /home
  currentStat    *SectionStat     // current set of statistics
  statHistory    []*SectionStat
}

// StartNewStats will create a new set of statistics
// for a section, and include it in that section's
// history
func (sectionData *SectionData) StartNewStats() {
  newStat := new(SectionStat)
  newStat.startTime =  time.Now()
  newStat.methodsCounts = make(map[string]int, 0)
  newStat.contentLengths = make([]int, 0)
  newStat.statusCounts = make(map[string]int, 0)
  newStat.hostCounts = make(map[string]int, 0)
  newStat.hitCount = 0

  sectionData.currentStat = newStat
  sectionData.statHistory = append(sectionData.statHistory, sectionData.currentStat)
}

var allSectionData map[string]*SectionData // all section data, by section name

func startStats(dataChan chan LineData) {
  allSectionData = make(map[string]*SectionData)

  for {
    select {
    case lineData := <-dataChan:
      updateStats(lineData)
    }
  }
}

func updateStats(data LineData) {
  sectionData := dataForSection(data.SectionStr)
  currentStats := sectionData.currentStat

  currentStats.hitCount++
  currentStats.methodsCounts[data.RequestMethod]++
  currentStats.statusCounts[data.Status]++
  currentStats.hostCounts[data.RemoteHost]++
  currentStats.contentLengths = append(currentStats.contentLengths, data.ContentLen)
}

func dataForSection(sectionName string) *SectionData {
  sectionData, present := allSectionData[sectionName]
  if !present {
    sectionData = new(SectionData)
    sectionData.name = sectionName
    sectionData.statHistory = make([]*SectionStat, 0)
    sectionData.StartNewStats()
    allSectionData[sectionName] = sectionData
  }

  return sectionData
}

func archiveStatData() {
  for index := range allSectionData {
    sectionData := allSectionData[index]
    sectionData.StartNewStats()
  }
}
