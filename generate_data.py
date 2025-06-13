import random
import csv

# 生成するアイテムの数
NUM_ITEMS = 1000

# 価値(value)のランダム範囲
MIN_VALUE = 10
MAX_VALUE = 1000

# 重さ(weight)のランダム範囲
MIN_WEIGHT = 5
MAX_WEIGHT = 200

# 出力するファイル名
OUTPUT_FILE = "knapsack_data.csv"


def generate_knapsack_data(num_items, min_val, max_val, min_w, max_w):
    """ナップサック問題のテストデータを生成する"""
    header = ["index", "value", "weight"]
    data = [header]

    for i in range(1, num_items + 1):
        value = random.randint(min_val, max_val)
        weight = random.randint(min_w, max_w)
        data.append([i, value, weight])

    return data


def save_to_csv(data, filename):
    """データをCSVファイルに保存する"""
    try:
        with open(filename, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerows(data)
        print(f"テストデータを '{filename}' に正常に保存しました。")
    except IOError as e:
        print(f"エラー: ファイルの書き込みに失敗しました。 {e}")


if __name__ == '__main__':
    # データの生成
    item_data = generate_knapsack_data(NUM_ITEMS, MIN_VALUE, MAX_VALUE, MIN_WEIGHT, MAX_WEIGHT)

    # ファイルへの保存
    save_to_csv(item_data, OUTPUT_FILE)
