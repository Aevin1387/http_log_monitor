package main

import (
  "time"
)

type SectionStat struct {
  hitCount       int              // number of hits
  methodsCounts  map[string] int  // number of hits per method
  contentLengths []int            // size of each hit
  statusCounts   map[string] int  // number of hits per status
  hostCounts     map[string] int  // number of hits from each host
  startTime      time.Time        // time when statistics started
}

type SectionData struct {
  name           string           // section name, ex /home
  currentStat    *SectionStat     // current set of statistics
  statHistory    []*SectionStat
  openAlert      *Alert           // an open alert, if any
  alerts         []*Alert         // all alerts on this section
}

const (
  tenSeconds = time.Second * 10
)

func (sectionData *SectionData) CurrentStats() *SectionStat {
  if time.Now().Sub(sectionData.currentStat.startTime) > tenSeconds {
    sectionData.StartNewStats()
  }

  return sectionData.currentStat
}

func (sectionData *SectionData) StartNewStats() {
  newStat := SectionStat{ startTime: time.Now() }
  newStat.methodsCounts = make(map[string]int, 0)
  newStat.contentLengths = make([]int, 0)
  newStat.statusCounts = make(map[string]int, 0)
  newStat.hostCounts = make(map[string]int, 0)

  sectionData.currentStat = &newStat
  sectionData.statHistory = append(sectionData.statHistory, sectionData.currentStat)
}


type Alert struct {
  section       *SectionStat      // section for alert
  hits          int               // number of hits causing alert
  alertedAt     time.Time         // time alert started
  resolvedAt    time.Time         // time alert resolved
}

var allSectionData map[string]SectionData // all section data, by section name

func startStats(data_chan chan LineData) {
  allSectionData = make(map[string]SectionData)
  for {
    select {
    case lineData := <-data_chan:

      updateStats(lineData)

    }
  }
}

func updateStats(data LineData) {
  sectionData := dataForSection(data.SectionStr)
  currentStats := sectionData.CurrentStats()
  currentStats.hitCount++
  currentStats.methodsCounts[data.RequestMethod]++
  currentStats.statusCounts[data.Status]++
  currentStats.hostCounts[data.RemoteHost]++
  currentStats.contentLengths = append(currentStats.contentLengths, data.ContentLen)
}

func dataForSection(sectionName string) SectionData {
  sectionData, present := allSectionData[sectionName]
  if !present {
    sectionData = *new(SectionData)
    sectionData.statHistory = make([]*SectionStat, 0)
    sectionData.alerts = make([]*Alert, 0)
    sectionData.StartNewStats()
    allSectionData[sectionName] = sectionData
  }

  return sectionData
}





