# ryot

最初期Git（commit-01）に寄せた、互換不要のミニ実装です。

- オブジェクトDB: `.dircache/objects/`（zlib圧縮したオブジェクトを保存）
- index（キャッシュ）: `.dircache/index`（自前バイナリ）

## links
- git もっとも古いcommitlog  
  https://git.kernel.org/pub/scm/git/git.git/log/?ofs=79200
- git 一番最初のcommit  
  https://git.kernel.org/pub/scm/git/git.git/commit/?id=e83c5163316f89bfbde7d9ab23ca2e25604af290

## build
リポジトリルートで：

make build
