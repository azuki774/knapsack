package scenario

import (
	"encoding/csv"
	"fmt"
	"knapsack/model"
	"log"
	"os"
	"sort"
	"strconv"
)

type Scenario struct {
	Treasures   []model.Treasure
	WeightLimit float64
}

type Strategy struct {
	ChooseIndex []bool
	SumValue    float64
	SumWeight   float64
	Score       float64
}

type Gene struct {
	Strategy  Strategy
	Generation int // 世代数
}

func NewStrategy(ts []model.Treasure) *Strategy {
	strategy := Strategy{
		ChooseIndex: make([]bool, len(ts)),
	}
	return &strategy
}

func (s *Scenario) greedy() {
	treasures := make([]model.Treasure, len(s.Treasures))
	copy(treasures, s.Treasures)
	// コスパ良い順に安定ソート
	sort.SliceStable(treasures, func(i, j int) bool {
		return treasures[i].Value/treasures[i].Weight > treasures[j].Value/treasures[j].Weight
	})

	strategy := NewStrategy(treasures)

	for _, treasure := range treasures {
		// 取って超えない場合取る
		if strategy.SumWeight+treasure.Weight <= s.WeightLimit {
			strategy.ChooseIndex[treasure.Index] = true
			strategy.SumValue += treasure.Value
			strategy.SumWeight += treasure.Weight
			strategy.Score += treasure.Value
		}
	}
	fmt.Println(strategy)
}

func (s *Scenario) Load() {
	// knapsack_data.csv を読み込む
	f, err := os.Open("knapsack_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll() // csvを一度に全て読み込む
	if err != nil {
		log.Fatal(err)
	}

	// [][]stringなのでループする
	for i, v := range rows {
		// ヘッダ行はスキップ
		if i == 0 {
			continue
		}
		// index, value, weight
		index, _ := strconv.Atoi(v[0])
		value, _ := strconv.Atoi(v[1])
		weight, _ := strconv.Atoi(v[2])
		treasure := model.Treasure{
			Index:  index,
			Value:  float64(value),
			Weight: float64(weight),
		}
		s.Treasures = append(s.Treasures, treasure)
	}
}

func (s *Scenario) GA() {
	// 初期集団の生成

	// 評価

	// 選択

	// 交叉

	// 世代交代
}

func (s *Scenario) Start() {
	s.Load()
	s.greedy()

}
