# gomockとは

Goのモックを作るためのフレームワークです。
モックを自動生成するツール`mockgen`とテスティングライブラリの`gomock`の2つを含めてgomockと呼ばれています。

golang/mock
https://github.com/golang/mock

# mockgenのインストール

```bash
go get github.com/golang/mock/mockgen
```

# gomockを使ったテストのやり方

以下のファイルを用意します。
```
.
├── application
│   └── task_application_service.go
├── domain
│   ├── task.go
│   └── task_repository.go
├── go.mod
└── go.sum
```

```go:task_application_service.go
package application

import (
	"context"
	"errors"
	"gomock_test/domain"
)

type TaskApplicationService struct {
	taskRepository domain.TaskRepository
}

func NewTaskApplicationService(taskRepository domain.TaskRepository) *TaskApplicationService {
	return &TaskApplicationService{taskRepository: taskRepository}
}

func (s *TaskApplicationService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	t, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		errors.New("get task failed.")
		return nil, err
	}

	return t, nil
}
```

```go:task.go
package domain

type Task struct {
	ID string
}
```

```go:task_repository.go
package domain

import "context"

type TaskRepository interface {
	Get(ctx context.Context, id string) (*Task, error)
}
```

```:go.mod
module gomock_test

go 1.12
```

## モックを作成する

interfaceが含まれいているソースコードを指定するとmockgenがそれに適したモックを自動生成します。

```
mockgen -source domain/task_repository.go -destination domain/mock_task_repository.go -package domain -self_package gomock_test/domain
```
### よく使うオプション
- `-source`: モック化の対象となるinterfaceが含まれるソースコードのファイル名
- `-destination`: モックファイルの出力先ファイル名
- `-package`: モックファイルのパッケージ名。指定しない場合は`mock_[カレントディレクトリ名]`となる。
- `-self_package`: モックが自身のパッケージを参照するのを避けるために使用する。

## モックを使ってテストをする

`gomock`を`go get`します。

```bash
go get github.com/golang/mock/gomock
```

`task_application_service`のテストファイル`task_application_service_test.go`を作ります。

```go
package application

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"gomock_test/domain"
	"testing"
)

func TestTaskApplicationService_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)

	taskRepository := domain.NewMockTaskRepository(ctrl)

	s := NewTaskApplicationService(taskRepository)

	t.Run("ok", func(t *testing.T) {
		task := &domain.Task{ID: "TEST_TASK"}
		taskRepository.EXPECT().Get(gomock.Any(), "TEST_TASK").Return(task, nil)

		ctx := context.Background()
		id := "TEST_TASK"

		task, err := s.GetTask(ctx, id)

		if err != nil {
			t.Fatal("expected no error")
		}
	})

	t.Run("fail", func(t *testing.T) {
		taskRepository.EXPECT().Get(gomock.Any(), "TEST_TASK").Return(nil, errors.New("something failed"))

		ctx := context.Background()
		id := "TEST_TASK"

		task, err := s.GetTask(ctx, id)

		if task != nil || err == nil {
			t.Fatal("expected error")
		}
	})
}
```

モックを使うにはまずコントローラを作成します。
```go
ctrl := gomock.NewController(t)
```
作成したコントローラを使ってモックをNewします。
```go
taskRepository := domain.NewMockTaskRepository(ctrl)
```

`TaskApplicationService`は`TaskRepositry`interfaceを必要としているので、作成したモックを引数にしてNewします。

```go
s := NewTaskApplicationService(taskRepository)
```

次にモックの挙動を設定します。
`EXPECT()`を呼ぶことで、メソッドの呼び出されたかどうかをテストします。
interfaceで定義したメソッド（今回の場合は`Get()`）の引数に期待される引数を入力します。
`Return()`の引数にmockの戻り値を入力します。
```go
taskRepository.EXPECT().Get(gomock.Any(), "TEST_TASK").Return(task, nil)
```
最後にテストしたいメソッドを呼び出して期待されている挙動かどうかをテストします。
```go
task, err := s.GetTask(ctx, id)

if task != nil || err == nil {
	t.Fatal("expected error")
}
```

mockで戻り値を固定化することにより、実際の実行結果に左右されず、エラーハンドリングをテストすることが容易になります。

# 参考
https://qiita.com/yuina1056/items/faf1aca0c052dfe802fe