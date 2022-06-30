# go-todo

書籍「Go言語ハンズオン」Section6-4から学習がてら機能を追加していく形で開発。
とりあえず動くものを作成しただけなので、リファクタ未実施の状態。

## シーケンス図

単純にGETしているだけのものは省略

### マイページ `/home`

``` mermaid
sequenceDiagram
    autonumber
    alt memo
    Client->>+Server: メモ作成
    Note over Server: "Memo"テーブルでinsert
    Server-->>-Client: 作成完了後、`/home`にリダイレクト
    else list
    Client->>+Server: リスト作成
    Note over Server: "List"テーブルでinsert
    Server-->>-Client: 作成完了後、`/home`にリダイレクト
    end
```

### 修正 `POST /edit`

``` mermaid
sequenceDiagram
    autonumber
    Client->>+Server: 修正するtitleとmessageを送信
    Note over Server: "Memo"テーブルで対象idにupdate
    Server-->>-Client: 完了後、`/detail`にリダイレクト
```

### 削除 `POST /delte`

``` mermaid
sequenceDiagram
    autonumber
    Client->>+Server: POST
    Note over Server: "Memo"テーブルで対象idにdelete
    Server-->>-Client: 完了後、`/`にリダイレクト
```


## ER図

簡素に作る。
DBはSQLite3。PostgreSQLのようなtimestamp型は無いため、text型で代用。

``` mermaid
erDiagram

users {
    id integer
    created_at text
    updated_at text
    deleted_at text
    account text
    name text
    password text
    message text
}

comments {
    id integer
    created_at text
    updated_at text
    deleted_at text
    user_id integer
    memo_id integer
    message text
}

memos {
    id integer
    created_at text
    updated_at text
    deleted_at text
    address text
    title text
    message text
    user_id integer
    list_id integer
}

lists {
    id integer
    created_at text
    updated_at text
    deleted_at text
    user_id integer
    name text
    message text
}

users ||--o{ comments: ""
users ||--o{ memos: ""
users ||--o{ lists: ""
memos ||--o{ comments: ""
memos }o--o{ lists: ""
```


## ToDo

- [x] アカウント新規作成追加
  - [x] ログインページにアカウント作成ボタンを追加
  - [x] 新規アカウント作成ページ作成
  - [x] 作成ボタンを押したら"Users"テーブルにinsert、作成されるようにする
  - [x] 作成後の移行先は`/`(トップページ)にする
- [x] マイページで表示するアカウント情報を増やす
  - [x] メールアドレス
  - [x] ニックネーム
- [x] メモを作成するとタイトル部分がURLから取った文字列になるので、タイトルを設定できるようにする
- [x] トップページに表示されるメモ・リストの表示を、自分が投稿したもののみに修正する
- [x] メモ編集機能を追加
  - [x] 編集
  - [x] 削除
- [ ] テスト作る