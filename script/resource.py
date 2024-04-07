import json, os, sys, shutil, time, threading, base64
from concurrent.futures import ThreadPoolExecutor
from io import BytesIO
from queue import Queue
from UnityPy import Environment
from UnityPy.enums import ClassIDType
from zipfile import ZipFile

queue_out = Queue()
queue_in = Queue()
classes = ClassIDType.TextAsset, ClassIDType.Sprite, ClassIDType.AudioClip

class ByteReader:
    def __init__(self, data):
        self.data = data
        self.position = 0
    def readInt(self):
        self.position += 4
        return self.data[self.position - 4] ^ self.data[self.position - 3] << 8 ^ self.data[self.position - 2] << 16

def io():
    while True:
        item = queue_in.get()
        if item == None:
            break
        elif type(item) == list:
            env = Environment()
            for i in range(1, len(item)):
                env.load_file(item[0].read("assets/aa/Android/%s" % item[i][1]), name=item[i][0])
            queue_out.put(env)
            del env
        else:
            path, resource = item
            if type(resource) == BytesIO:
                with resource:
                    with open(path, "wb") as f:
                        f.write(resource.getbuffer())
            else:
                with open(path, "wb") as f:
                    f.write(resource)

def save_image(path, image):
    bytesIO = BytesIO()
    image.save(bytesIO, "png")
    queue_in.put((path, bytesIO))
def save(key, entry):
    obj = entry.get_filtered_objects(classes)
    obj = next(obj).read()
    if key[-25:] == ".0/IllustrationLowRes.png":
        key = key[:-25]
        pool.submit(save_image, "IllustrationLowRes/%s.png" % key, obj.image)

def run(path):
    with ZipFile(path) as apk:
        with apk.open("assets/aa/catalog.json") as f:
            data = json.load(f)

    for directory in ["IllustrationLowRes"]:
        shutil.rmtree(directory, True)
        os.mkdir(directory)

    key = base64.b64decode(data["m_KeyDataString"])
    bucket = base64.b64decode(data["m_BucketDataString"])
    entry = base64.b64decode(data["m_EntryDataString"])
    table, reader= [], ByteReader(bucket)

    for _ in range(reader.readInt()):
        key_position = reader.readInt()
        key_type = key[key_position]
        key_position += 1
        if key_type == 0:
            length = key[key_position]
            key_position += 4
            key_value = key[key_position:key_position + length].decode()
        elif key_type == 1:
            length = key[key_position]
            key_position += 4
            key_value = key[key_position:key_position + length].decode("utf16")
        elif key_type == 4:
            key_value = key[key_position]
        else:
            raise BaseException(key_position, key_type)
        for i in range(reader.readInt()):
            entry_position = reader.readInt()
            entry_value = entry[4 + 28 * entry_position:4 + 28 * entry_position + 28]
            entry_value = entry_value[8] ^ entry_value[9] << 8
        table.append([key_value, entry_value])
    for i in range(len(table)):
        if table[i][1] != 65535:
            table[i][1] = table[table[i][1]][0]
    for i in range(len(table) - 1, -1, -1):
        if type(table[i][0]) == int or table[i][0][:15] == "Assets/Tracks/#" or table[i][0][:14] != "Assets/Tracks/" and \
                table[i][0][:7] != "avatar.":
            del table[i]
        elif table[i][0][:14] == "Assets/Tracks/":
            table[i][0] = table[i][0][14:]

    thread = threading.Thread(target=io)
    thread.start()
    ti = time.time()
    global pool
    with ThreadPoolExecutor(6) as pool:
        with ZipFile(path) as apk:
            size, l = 0, [apk]
            for key, entry in table:
                l.append((key, entry))
                info = apk.getinfo("assets/aa/Android/%s" % entry)
                size += info.file_size
                if size > 32 * 1024 * 1024:
                    queue_in.put(l)
                    env = queue_out.get()
                    for ikey, ientry in env.files.items():
                        save(ikey, ientry)
                    size = 0
                    del env
                    l = [apk]
            queue_in.put(l)
            env = queue_out.get()
            for ikey, ientry in env.files.items():
                save(ikey, ientry)
    queue_in.put(None)
    thread.join()
    print("%f秒" % round(time.time() - ti, 4))

if __name__ == "__main__":
    run(sys.argv[1])
