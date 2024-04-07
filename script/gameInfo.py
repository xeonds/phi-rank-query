import struct
import sys
from UnityPy import Environment
import zipfile
import csv

SONG_BASE_SCHEMA = {
    "songId": str, "songKey": str, "songName": str, "songTitle": str, "difficulty": [float],
    "illustrator": str, "charter": [str], "composer": str, "levels": [str], "previewTimeFrom": float, "previewTimeTo": float,
    "unlockList": {"unlockType": int, "unlockInfo": [str]}, "levelMods": {"n": [str]}
}

class ByteReader:
    def __init__(self, data: bytes):
        self.data = data
        self.position = 0
        self.d = {int: self.readInt, float: self.readFloat, str: self.readString}

    # 4字节读取数据(int)
    def readInt(self):
        self.position += 4
        return self.data[self.position - 4] ^ self.data[self.position - 3] << 8

    # 4字节读取数据(float)
    def readFloat(self):
        self.position += 4
        return struct.unpack("f", self.data[self.position - 4:self.position])[0]

    # 读取字符串
    def readString(self):
        length = self.readInt() # 读取第一个字节获取当前字符串长度
        result = self.data[self.position:self.position + length].decode()
        self.position += length // 4 * 4
        if length % 4 != 0:
            self.position += 4
        return result

    def skipString(self): # 略过字符串
        length = self.readInt()
        self.position += length // 4 * 4
        if length % 4 != 0:
            self.position += 4

    def readSchema(self, schema: dict): # 通过SONG_BASE_SCHEMA中的字典来获取数据类型
        result = []
        for x in range(self.readInt()):
            item = {}
            for key, value in schema.items():
                if value in (int, str, float):
                    item[key] = self.d[value]()
                elif type(value) == list:
                    l = []
                    for i in range(self.readInt()):
                        l.append(self.d[value[0]]())
                    item[key] = l
                elif type(value) == tuple:
                    for t in value:
                        self.d[t]()
                elif type(value) == dict:
                    item[key] = self.readSchema(value)
                else:
                    raise Exception("无")
            result.append(item)
        return result

def run(path):
    env = Environment()
    with zipfile.ZipFile(path) as apk:
        with apk.open("assets/bin/Data/globalgamemanagers.assets") as f:
            env.load_file(f.read(), name="assets/bin/Data/globalgamemanagers.assets")
        with apk.open("assets/bin/Data/level0") as f:
            env.load_file(f.read())
    for obj in env.objects:
        if obj.type.name != "MonoBehaviour": continue
        data = obj.read()
        if data.m_Script.get_obj().read().name == "GameInformation":
            information = data.raw_data.tobytes()

    reader = ByteReader(information)
    reader.position = information.index(b"\x16\x00\x00\x00Glaciaxion.SunsetRay.0\x00\x00\n") - 4
    difficulty = []
    table = []
    for i in range(3):
        for item in reader.readSchema(SONG_BASE_SCHEMA):
            item["songId"] = item["songId"][:-2]
            if len(item["levels"]) == 5:
                item["difficulty"].pop()
                item["charter"].pop()
            if item["difficulty"][-1] == 0:
                item["difficulty"].pop()
                item["charter"].pop()
            for i in range(len(item["difficulty"])):
                item["difficulty"][i] = round(item["difficulty"][i], 1)
            difficulty.append([item["songId"]] + item["difficulty"])
            table.append([item["songId"], item["songName"], item["composer"], item["illustrator"], *item["charter"]])
            
    reader.readSchema(SONG_BASE_SCHEMA)
    print("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-")

    with open("difficulty.csv", "w", encoding="utf8", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(["id", "EZ", "HD", "IN", "AT"])
        for item in difficulty:
            while len(item) < 5: item.append("")
        writer.writerows(difficulty)
    print("difficulty write completed")

    with open("info.csv", "w", encoding="utf8", newline='') as f:
        writer = csv.writer(f)
        writer.writerow(["id", "song", "composer", "illustrator", "EZ", "HD", "IN", "AT"])
        for item in table:
            while len(item) < 8: item.append("")
        writer.writerows(table)
    print("info write completed")
    print("done.")

if __name__ == "__main__":
    run(sys.argv[1])
