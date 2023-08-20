"""
生成分词器的测试数据
"""
import json
import os
import token
import secrets
import shutil
import sys
import time
import tokenize


def random_filename() -> str:
    return secrets.token_hex(8) + '_' + str(int(time.time() * 1000))


def dump_tokens(dst_dir: str, filepath: str) -> None:
    with open(filepath, 'rb') as f:
        tokens = list(tokenize.tokenize(f.__next__))

    data_filepath = os.path.join(dst_dir, f'{random_filename()}.json')
    json_tokens = [{'type': tk.type, 'literal': tk.string, 'start': list(tk.start), 'end': list(tk.end)} for tk in tokens]
    dump_data = {
        'filepath': filepath,
        'tokens': json_tokens,
    }
    with open(data_filepath, 'w') as f:
        json.dump(dump_data, f, ensure_ascii=False, indent=4)


def main() -> None:
    if len(sys.argv) != 2:
        print(f"Usage: python {sys.argv[0]} <python filepath or directory>")
        return
    testdir = "/tmp/checkpy3lexer2023"
    if os.path.exists(testdir):
        shutil.rmtree(testdir)
    os.makedirs(str(testdir), exist_ok=True)

    for root, dirs, files in os.walk(sys.argv[1]):
        for filename in files:
            if not filename.endswith('.py'):
                continue
            print(filename)
            filepath = os.path.join(root, filename)
            filepath = os.path.abspath(filepath)
            dump_tokens(str(testdir), filepath)


if __name__ == '__main__':
    main()
