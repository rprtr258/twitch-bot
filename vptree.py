def hammingDistance(str1, str2):
    result = 0
    n = min(str1.length(), str2.length())
    result += max(str1.length(), str2.length()) - n
    for i in range(n):
        result += (str1[i] != str2[i])
    return result

def editDistance(str1, str2):
    m = len(str1)
    n = len(str2)
    if m == 0:
        return n
    if n == 0:
        return m
    matrix = [[0 for _ in range(n + 1)] for i in range(m + 1)]
    for i in range(1, m + 1):
        matrix[i][0] = i
    for i in range(1, n + 1):
        matrix[0][i] = i
    for i in range(m + 1):
        for j in range(1, n + 1):
            cost = 0 if str1[i - 1] == str2[j - 1] else 1
            above_cell = matrix[i - 1][j]
            left_cell = matrix[i][j - 1]
            diagonal_cell = matrix[i - 1][j - 1]
            matrix[i][j] = min(above_cell + 1, left_cell + 1, diagonal_cell + cost)
    return matrix[m][n]

def editIgnoreCaseDistance(str1, str2):
    return editDistance(str1.lower(), str2.lower())

class Node:
    def __init__(self, list, strDist):
        self.dist = strDist
        if len(list) == 1:
            self.radius = 0
            self.data = list[0]
            self.outer = None
            self.inner = None
            return
        self.data = list[0]
        list = list[1:]
        self.radius = self.dist(self.data, list[-1]) / 2
        inside = [x for x in list if self.dist(x, self.data) <= self.radius]
        outside = [x for x in list if self.dist(x, self.data) > self.radius]
        if inside == []:
            self.inner = None
        else:
            self.inner = Node(inside, self.dist)
        if outside == []:
            self.outer = None
        else:
            self.outer = Node(outside, self.dist)

    def findNearest(self, str, prec):
        d = self.dist(self.data, str)
        result = []
        if d <= prec:
            result.append(self.data)
        if d + prec >= self.radius and self.outer is not None:
            add = self.outer.findNearest(str, prec)
            result += add
        if d - prec <= self.radius and self.inner is not None:
            add = self.inner.findNearest(str, prec)
            result += add
        return result

class VPTree:
    def __init__(self, list, strDist):
        self.dist = strDist
        self.root = Node(list, strDist)

    def findNearest(self, str, prec):
        return self.root.findNearest(str, prec)

    def dist(self, a, b):
        return self.dist(a, b)
    
    def findSimilar(self, str):
        result = self.findNearest(str, len(str))
        result = sorted(result, key=lambda a: self.dist(a, str))
        return result
