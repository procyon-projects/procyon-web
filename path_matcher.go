package web

import (
	"regexp"
	"strings"
)

type PathMatcher interface {
	MatchPath(path string, pattern string) bool
	GetUriVariables(path string, pattern string) map[string]string
}

/*
 * SimplePatchMatcher was developed by looking into AntPathMatcher in Spring Framework
 * and AntPathMatcher's code was used.
 */

type SimplePathMatcher struct {
	pathSeparator string
	globalPattern *regexp.Regexp
}

func NewSimplePathMatcher() SimplePathMatcher {
	globalPattern, _ := regexp.Compile("\\?|\\*|\\{((?:\\{[^/]+?}|[^/{}]|\\\\[{}])+?)}")
	return SimplePathMatcher{
		pathSeparator: "/",
		globalPattern: globalPattern,
	}
}

func (pathMatcher SimplePathMatcher) MatchPath(path string, pattern string) bool {
	if pattern == path {
		return true
	}
	return pathMatcher.match(path, pattern, nil)
}

func (pathMatcher SimplePathMatcher) GetUriVariables(path string, pattern string) map[string]string {
	variables := make(map[string]string, 0)
	pathMatcher.match(path, pattern, variables)
	return variables
}

func (pathMatcher SimplePathMatcher) match(path string, pattern string, variables map[string]string) bool {
	if strings.HasPrefix(path, pathMatcher.pathSeparator) && strings.HasPrefix(pattern, pathMatcher.pathSeparator) {
		patternDirs := strings.Split(pattern, pathMatcher.pathSeparator)
		if !pathMatcher.isCandidate(patternDirs, path) {
			return false
		}

		pathDirs := strings.Split(path, pathMatcher.pathSeparator)
		patternStartIndex := 0
		patternEndIndex := len(patternDirs) - 1
		pathStartIndex := 0
		pathEndIndex := len(pathDirs) - 1
		patternDir := ""
		for ; patternStartIndex <= patternEndIndex && pathStartIndex <= pathEndIndex; pathStartIndex++ {
			patternDir = patternDirs[patternStartIndex]
			if "**" == patternDir {
				break
			}
			if !pathMatcher.matchWithRegex(pathDirs[pathStartIndex], patternDir, variables) {
				return false
			}
			patternStartIndex++
		}

		if pathStartIndex > pathEndIndex {
			if patternStartIndex > patternEndIndex {
				return true
			} else {
				patternTempIndex := patternStartIndex
				for ; patternTempIndex <= patternEndIndex; patternTempIndex++ {
					if patternDirs[patternTempIndex] != "**" {
						return false
					}
				}
				return true
			}
		} else if patternStartIndex > patternEndIndex {
			return false
		} else {
			for patternStartIndex <= patternEndIndex && pathStartIndex <= pathEndIndex {
				patternDir = patternDirs[patternEndIndex]
				if patternDir == "**" {
					break
				}
				if !pathMatcher.matchWithRegex(pathDirs[pathEndIndex], patternDir, variables) {
					return false
				}
				pathEndIndex--
				patternEndIndex--
			}
			if patternStartIndex > patternEndIndex {
				patternTempIndex := patternStartIndex
				for ; patternTempIndex <= patternEndIndex; patternTempIndex++ {
					if patternDirs[patternTempIndex] != "**" {
						return false
					}
				}
				return true
			} else {
				patternIndexTemp := 0
				for patternStartIndex != patternEndIndex && pathStartIndex <= pathEndIndex {
					patternIndexTemp = -1

					patternLength := patternStartIndex + 1
					for ; patternLength <= patternEndIndex; patternLength++ {
						if patternDirs[patternLength] == "**" {
							patternIndexTemp = patternLength
							break
						}
					}

					if patternIndexTemp == patternStartIndex+1 {
						patternStartIndex++
					} else {
						patternLength = patternIndexTemp - patternStartIndex - 1
						strLength := pathEndIndex - patternStartIndex + 1
						foundIndex := -1

					searchPattern:
						for i := 0; i <= strLength-patternLength; {
							for j := 0; j < patternLength; j++ {
								subPat := patternDirs[patternStartIndex+j+1]
								subStr := pathDirs[pathStartIndex+i+j]
								if !pathMatcher.matchWithRegex(subPat, subStr, variables) {
									i++
									continue searchPattern
								}
							}

							foundIndex = patternStartIndex + i
							break
						}

						if foundIndex == -1 {
							return false
						}

						patternStartIndex = patternIndexTemp
						pathStartIndex = foundIndex + patternLength
					}
				}

				for patternIndexTemp = patternStartIndex; patternIndexTemp <= patternEndIndex; patternIndexTemp++ {
					if patternDirs[patternIndexTemp] != "**" {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func (pathMatcher SimplePathMatcher) isCandidate(patternDirs []string, path string) bool {
	position := 0
	for _, pattern := range patternDirs {
		skipped := pathMatcher.skipSeparator(path, position)
		position += skipped
		skipped = pathMatcher.skipPathSegment(path, position, pattern)
		if skipped < len(pattern) {
			return skipped > 0 || len(pattern) > 0 && pathMatcher.isWildcardChar([]rune(pattern)[0])
		}
		position += skipped
	}
	return true
}

func (pathMatcher SimplePathMatcher) skipSeparator(path string, position int) int {
	skipped := 0
	for strings.HasPrefix(path[position+skipped:], pathMatcher.pathSeparator) {
		skipped += len(pathMatcher.pathSeparator)
	}
	return skipped
}

func (pathMatcher SimplePathMatcher) skipPathSegment(path string, position int, patternDir string) int {
	skipped := 0
	for _, character := range patternDir {
		if pathMatcher.isWildcardChar(character) {
			return skipped
		}
		currentPosition := position + skipped
		if currentPosition >= len(path) {
			return 0
		}
		if character == []rune(path)[currentPosition] {
			skipped++
		}
	}
	return skipped
}

func (pathMatcher SimplePathMatcher) isWildcardChar(character int32) bool {
	if character == '?' || character == '*' || character == '{' {
		return true
	}
	return false
}

func (pathMatcher SimplePathMatcher) matchWithRegex(str string, pattern string, variables map[string]string) bool {
	regex, err := pathMatcher.getRegexPattern(pattern)
	if err != nil {
		return false
	}
	result := regex.FindStringSubmatch(str)
	if result != nil && variables != nil {
		for index, expName := range regex.SubexpNames() {
			if len(expName) != 0 {
				variables[expName] = result[index]
			}
		}
	}
	return true
}

func (pathMatcher SimplePathMatcher) getRegexPattern(pattern string) (*regexp.Regexp, error) {
	result := pathMatcher.globalPattern.FindAllStringSubmatch(pattern, -1)
	resultIndex := pathMatcher.globalPattern.FindAllStringSubmatchIndex(pattern, -1)

	startIndex := 0
	buildPattern := ""
	if result != nil && resultIndex != nil {
		keyword := result[0][0]
		keywordStartIndex := resultIndex[0][0]
		if startIndex != keywordStartIndex {
			buildPattern = buildPattern + pattern[startIndex:keywordStartIndex]
		}
		startIndex = keywordStartIndex + len(keyword)
		if keyword == "?" {
			buildPattern = buildPattern + "."
		} else if keyword == "*" {
			buildPattern = buildPattern + ".*"
		} else if strings.HasPrefix(keyword, "{") && strings.HasSuffix(keyword, "}") {
			colonIndex := strings.Index(keyword, ":")
			if colonIndex == -1 {
				variableName := keyword[1 : len(keyword)-1]
				buildPattern = buildPattern + "(?P<" + variableName + ">.*)"
			} else {
				variablePattern := keyword[colonIndex+1 : len(keyword)-1]
				variableName := keyword[1:colonIndex]
				buildPattern = buildPattern + "(?P<" + variableName + ">" + variablePattern + ")"
			}
		}
		if startIndex != len(pattern) {
			buildPattern = buildPattern + pattern[startIndex:]
		}
	} else {
		buildPattern = pattern
	}
	return regexp.Compile("^" + buildPattern + "$")
}
