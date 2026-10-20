package services

import (
	"encoding/json"
	"log"

	"devcode/backend/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedProblems(db *gorm.DB) error {
	var count int64
	db.Model(&models.Problem{}).Count(&count)
	if count > 0 {
		return nil
	}

	problems := []struct {
		Problem   models.Problem
		TestCases []models.TestCase
	}{
		{
			Problem: models.Problem{
				Slug:        "two-sum",
				Title:       "Two Sum",
				Difficulty:  "easy",
				Description: "Given an array of integers `nums` and an integer `target`, return indices of the two numbers such that they add up to `target`.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.\n\n### Input Format\nFirst line: space-separated integers for `nums`\nSecond line: integer `target`\n\n### Output Format\nSpace-separated zero-based indices in ascending order (e.g. `0 1`)",
				StarterCode: toJSON(map[string]string{
					"python":     "import sys\n\ndef two_sum(nums, target):\n    # Write your solution here\n    seen = {}\n    for i, num in enumerate(nums):\n        diff = target - num\n        if diff in seen:\n            return f\"{seen[diff]} {i}\"\n        seen[num] = i\n    return \"\"\n\nif __name__ == '__main__':\n    lines = sys.stdin.read().splitlines()\n    if lines:\n        nums = list(map(int, lines[0].split()))\n        target = int(lines[1])\n        print(two_sum(nums, target))\n",
					"javascript": "const fs = require('fs');\n\nfunction twoSum(nums, target) {\n    // Write your solution here\n    const seen = new Map();\n    for (let i = 0; i < nums.length; i++) {\n        const diff = target - nums[i];\n        if (seen.has(diff)) {\n            return `${seen.get(diff)} ${i}`;\n        }\n        seen.set(nums[i], i);\n    }\n    return '';\n}\n\nconst input = fs.readFileSync(0, 'utf-8').trim().split('\\n');\nif (input.length >= 2) {\n    const nums = input[0].trim().split(/\\s+/).map(Number);\n    const target = Number(input[1].trim());\n    console.log(twoSum(nums, target));\n}\n",
					"typescript": "const fs = require('fs');\n\nfunction twoSum(nums: number[], target: number): string {\n    const seen = new Map<number, number>();\n    for (let i = 0; i < nums.length; i++) {\n        const diff = target - nums[i];\n        if (seen.has(diff)) {\n            return `${seen.get(diff)} ${i}`;\n        }\n        seen.set(nums[i], i);\n    }\n    return '';\n}\n\nconst input = fs.readFileSync(0, 'utf-8').trim().split('\\n');\nif (input.length >= 2) {\n    const nums = input[0].trim().split(/\\s+/).map(Number);\n    const target = Number(input[1].trim());\n    console.log(twoSum(nums, target));\n}\n",
					"go":         "package main\n\nimport (\n\t\"bufio\"\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n\t\"strings\"\n)\n\nfunc main() {\n\tscanner := bufio.NewScanner(os.Stdin)\n\tif !scanner.Scan() {\n\t\treturn\n\t}\n\tnumsStr := strings.Fields(scanner.Text())\n\tif !scanner.Scan() {\n\t\treturn\n\t}\n\ttarget, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))\n\tseen := make(map[int]int)\n\tfor i, s := range numsStr {\n\t\tnum, _ := strconv.Atoi(s)\n\t\tdiff := target - num\n\t\tif prev, ok := seen[diff]; ok {\n\t\t\tfmt.Printf(\"%d %d\\n\", prev, i)\n\t\t\treturn\n\t\t}\n\t\tseen[num] = i\n\t}\n}\n",
				}),
				TimeLimitMs:   2000,
				MemoryLimitKb: 131072,
			},
			TestCases: []models.TestCase{
				{Input: "2 7 11 15\n9", ExpectedOutput: "0 1", IsHidden: false, SequenceOrder: 1},
				{Input: "3 2 4\n6", ExpectedOutput: "1 2", IsHidden: false, SequenceOrder: 2},
				{Input: "3 3\n6", ExpectedOutput: "0 1", IsHidden: false, SequenceOrder: 3},
				{Input: "1 5 8 12 19 23\n27", ExpectedOutput: "2 4", IsHidden: true, SequenceOrder: 4},
			},
		},
		{
			Problem: models.Problem{
				Slug:        "valid-palindrome",
				Title:       "Valid Palindrome",
				Difficulty:  "easy",
				Description: "A phrase is a **palindrome** if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward.\n\n### Input Format\nA single string `s`.\n\n### Output Format\nOutput `true` or `false`.",
				StarterCode: toJSON(map[string]string{
					"python":     "import sys, re\n\ndef is_palindrome(s):\n    cleaned = re.sub(r'[^a-zA-Z0-9]', '', s).lower()\n    return 'true' if cleaned == cleaned[::-1] else 'false'\n\nif __name__ == '__main__':\n    s = sys.stdin.read().rstrip('\\r\\n')\n    print(is_palindrome(s))\n",
					"javascript": "const fs = require('fs');\n\nfunction isPalindrome(s) {\n    const cleaned = s.replace(/[^a-zA-Z0-9]/g, '').toLowerCase();\n    const rev = cleaned.split('').reverse().join('');\n    return cleaned === rev ? 'true' : 'false';\n}\n\nconst input = fs.readFileSync(0, 'utf-8').replace(/[\\r\\n]+$/, '');\nconsole.log(isPalindrome(input));\n",
					"typescript": "const fs = require('fs');\n\nfunction isPalindrome(s: string): string {\n    const cleaned = s.replace(/[^a-zA-Z0-9]/g, '').toLowerCase();\n    const rev = cleaned.split('').reverse().join('');\n    return cleaned === rev ? 'true' : 'false';\n}\n\nconst input = fs.readFileSync(0, 'utf-8').replace(/[\\r\\n]+$/, '');\nconsole.log(isPalindrome(input));\n",
					"go":         "package main\n\nimport (\n\t\"bufio\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"unicode\"\n)\n\nfunc main() {\n\tscanner := bufio.NewScanner(os.Stdin)\n\tif scanner.Scan() {\n\t\ts := scanner.Text()\n\t\tvar runes []rune\n\t\tfor _, r := range s {\n\t\t\tif unicode.IsLetter(r) || unicode.IsDigit(r) {\n\t\t\t\trunes = append(runes, unicode.ToLower(r))\n\t\t\t}\n\t\t}\n\t\tn := len(runes)\n\t\tfor i := 0; i < n/2; i++ {\n\t\t\tif runes[i] != runes[n-1-i] {\n\t\t\t\tfmt.Println(\"false\")\n\t\t\t\treturn\n\t\t\t}\n\t\t}\n\t\tfmt.Println(\"true\")\n\t}\n}\n",
				}),
				TimeLimitMs:   2000,
				MemoryLimitKb: 131072,
			},
			TestCases: []models.TestCase{
				{Input: "A man, a plan, a canal: Panama", ExpectedOutput: "true", IsHidden: false, SequenceOrder: 1},
				{Input: "race a car", ExpectedOutput: "false", IsHidden: false, SequenceOrder: 2},
				{Input: " ", ExpectedOutput: "true", IsHidden: false, SequenceOrder: 3},
				{Input: "0P", ExpectedOutput: "false", IsHidden: true, SequenceOrder: 4},
			},
		},
		{
			Problem: models.Problem{
				Slug:        "reverse-words",
				Title:       "Reverse Words in a String",
				Difficulty:  "medium",
				Description: "Given an input string `s`, reverse the order of the **words**.\n\nA word is defined as a sequence of non-space characters. The words in `s` will be separated by at least one space.\n\nReturn a string of the words in reverse order concatenated by a single space, without leading or trailing spaces.\n\n### Input Format\nA single string `s`.\n\n### Output Format\nThe words reversed, separated by a single space.",
				StarterCode: toJSON(map[string]string{
					"python":     "import sys\n\ndef reverse_words(s):\n    words = s.strip().split()\n    return ' '.join(reversed(words))\n\nif __name__ == '__main__':\n    s = sys.stdin.read().rstrip('\\r\\n')\n    print(reverse_words(s))\n",
					"javascript": "const fs = require('fs');\n\nfunction reverseWords(s) {\n    return s.trim().split(/\\s+/).reverse().join(' ');\n}\n\nconst input = fs.readFileSync(0, 'utf-8').replace(/[\\r\\n]+$/, '');\nconsole.log(reverseWords(input));\n",
					"typescript": "const fs = require('fs');\n\nfunction reverseWords(s: string): string {\n    return s.trim().split(/\\s+/).reverse().join(' ');\n}\n\nconst input = fs.readFileSync(0, 'utf-8').replace(/[\\r\\n]+$/, '');\nconsole.log(reverseWords(input));\n",
					"go":         "package main\n\nimport (\n\t\"bufio\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n)\n\nfunc main() {\n\tscanner := bufio.NewScanner(os.Stdin)\n\tif scanner.Scan() {\n\t\twords := strings.Fields(scanner.Text())\n\t\tfor i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {\n\t\t\twords[i], words[j] = words[j], words[i]\n\t\t}\n\t\tfmt.Println(strings.Join(words, \" \"))\n\t}\n}\n",
				}),
				TimeLimitMs:   2000,
				MemoryLimitKb: 131072,
			},
			TestCases: []models.TestCase{
				{Input: "the sky is blue", ExpectedOutput: "blue is sky the", IsHidden: false, SequenceOrder: 1},
				{Input: "  hello world  ", ExpectedOutput: "world hello", IsHidden: false, SequenceOrder: 2},
				{Input: "a good   example", ExpectedOutput: "example good a", IsHidden: false, SequenceOrder: 3},
			},
		},
	}

	for _, item := range problems {
		p := item.Problem
		if err := db.Create(&p).Error; err != nil {
			log.Printf("Failed to seed problem %s: %v", p.Slug, err)
			continue
		}
		for _, tc := range item.TestCases {
			tc.ProblemID = p.ID
			db.Create(&tc)
		}
	}

	log.Println("✅ Initial algorithmic problems seeded")
	return nil
}

func toJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
