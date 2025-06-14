package scenario

import (
	"encoding/csv"
	"fmt"
	"knapsack/model"
	"log"
	"math/rand/v2"
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
	Strategy Strategy
}

type Population struct {
	Genes      []Gene
	Generation int // 世代数
}

type CrossoverPlan struct {
	Parent1 int
	Parent2 int
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

// roulette 選択して、親を選択
func (s *Scenario) rouletteSelect(population Population, totalScore uint64) int {
	// 選択
	pickV := rand.N(totalScore + 1) // 0 から totalScore

	// 集団の最初の個体から順に適応度を引いていく
	var currentSum float64 = 0
	for i, g := range population.Genes {
		currentSum += g.Strategy.Score
		// ランダムな値（矢）が現在の個体の範囲に入ったら、その個体を選択
		if float64(pickV) <= currentSum {
			return i
		}
	}

	// ごく稀な浮動小数点数の誤差などでループを抜けてしまった場合の安全策として、
	// 最後の個体の index を返す
	return len(population.Genes) - 1
}

// 一点交叉する
func singlePointCrossover(p1, p2 []bool) (c1 []bool, c2 []bool, crossover int) {
	l := len(p1)
	if len(p1) != len(p2) {
		fmt.Println("length error")
		panic("length error")
	}

	if l < 2 {
		return p1, p2, 0
	}

	crossoverPoint := rand.N(l-1) + 1 // 交叉は端では行わない .. 2 だったら index=1 と index=2 の間に cut を入れる
	c1 = make([]bool, l)
	c2 = make([]bool, l)

	// 子1 = 親1の前半 + 親2の後半
	copy(c1[:crossoverPoint], p1[:crossoverPoint])
	copy(c1[crossoverPoint:], p2[crossoverPoint:])

	// 子2 = 親2の前半 + 親1の後半
	copy(c2[:crossoverPoint], p2[:crossoverPoint])
	copy(c2[crossoverPoint:], p1[crossoverPoint:])

	return c1, c2, crossover
}

func (s *Scenario) GA() {
	var population Population
	population.Generation = 1 // 第一世代
	// 初期集団の生成
	for i := 0; i < model.PopulationSize; i++ {
		ns := NewStrategy(s.Treasures)
		for j := 0; j < len(s.Treasures); j++ {
			// i 番目の Geme の設定
			// j 番目の Treasure を取得するかどうか
			b := rand.N(2) // 0 or 1
			if b == 1 {
				ns.ChooseIndex[j] = true
				ns.SumValue += s.Treasures[j].Value
				ns.SumWeight += s.Treasures[j].Weight
				// Score は計算しない
			}
		}
		// 初期 popluration を追加
		population.Genes = append(population.Genes, Gene{
			Strategy: *ns,
		})
	}

	// 評価(score)を計算
	var totalScore uint64 = 0 // スコアは整数
	for i := range population.Genes {
		strategyPtr := &population.Genes[i].Strategy // for 内でスライスの値を変更するためにポインタ取得
		if strategyPtr.SumWeight > s.WeightLimit {
			strategyPtr.Score = 0 // 積載超過していたら0点
		} else {
			strategyPtr.Score = strategyPtr.SumValue
		}
		totalScore += uint64(strategyPtr.Score)
	}
	fmt.Printf("%+v\n", population)
	fmt.Printf("total_score: %v\n", totalScore)

	nextPopulation := Population{
		Genes:      make([]Gene, 0),
		Generation: population.Generation + 1,
	}

	for i := 0; i < len(s.Treasures)/2; i++ { // 選択・交叉の回数
		// 選択
		p1 := s.rouletteSelect(population, totalScore)
		p2 := s.rouletteSelect(population, totalScore)
		fmt.Println(p1, p2)
		// TODO: 同じ親を選ばないようにする

		// 交叉
		c1, c2, _ := singlePointCrossover(
			population.Genes[p1].Strategy.ChooseIndex,
			population.Genes[p2].Strategy.ChooseIndex,
		)

		nextPopulation.Genes = append(nextPopulation.Genes, Gene{
			Strategy: Strategy{
				ChooseIndex: c1,
			},
		})
		nextPopulation.Genes = append(nextPopulation.Genes, Gene{
			Strategy: Strategy{
				ChooseIndex: c2,
			},
		})

	}

	// 世代交代
}

func (s *Scenario) Start() {
	s.Load()
	s.greedy()
	fmt.Println("------")
	s.GA()
}
