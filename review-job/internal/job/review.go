package job

// 实现从kafka里面拉取数据的变更放到elasticsearch中

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"review-job/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"
)

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

// JobWorker自定义实现job结构体，是按transport.Server的接口
type JobWorker struct {
	kafkaReader *kafka.Reader
	esClient    *ESClient
	log         *log.Helper
}

type ESClient struct {
	*elasticsearch.TypedClient
	index string
}

// Msg定义kafkazzh
type Msg struct {
	Type     string `json:"type"`
	Database string `json:"database"`
	Table    string `json:"table"`
	IsDdl    bool `json:"isDdl"`
	Data     []map[string]interface{} `json:"data"`
}

func NewJobWorker(kafkaReader *kafka.Reader, esClient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafkaReader,
		esClient:    esClient,
		log:         log.NewHelper(logger),
	}
}

func NewKafkaReader(conf *conf.Kafka) *kafka.Reader {
	// 创建一个reader，指定GroupID，从 topic-A 消费消息
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: conf.GroupId, // 指定消费者组id
		Topic:   conf.Topic,
	})
}

func NewesClient(conf *conf.Elasticsearch) (*ESClient, error) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: conf.Addresses,
	}
	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	return &ESClient{
		TypedClient: client,
		index:       conf.Index,
	}, err
}

// Start kratos启动之后会调用的
// ctx是kratos框架启动的时候传入的ctx，是带有退出取消的
func (j *JobWorker) Start(ctx context.Context) error {
	j.log.Debug("JobWorker start")
	// 1. 从kafka获取mysql中数据变更的消息
	// 接收消息
	for {
		m, err := j.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) { // 这里是优雅的实现kratos停止后返回预期的结束
			return nil
		}
		if err != nil {
			j.log.Errorf("ReadMessage failed: %#v", err)
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 2. 将完整的评价数据写入ES
		msg := new(Msg)
		if err = json.Unmarshal(m.Value, msg); err != nil {
			j.log.Errorf("Unmarshal failed: %#v", err)
			continue
		}
		if msg.Type == "INSERT" {
			// 往ES中新增文档
			for idx := range msg.Data {
				j.indexDocument(msg.Data[idx])
			}

		} else {
			// 往ES中跟新文档
			for idx := range msg.Data {
				j.updateDocument(msg.Data[idx])
			}
		}
	}

	return nil
}

// Stop kratos结束之后会调用的
func (j *JobWorker) Stop(ctx context.Context) error {
	// 程序退出之前关闭Reader
	j.log.Debug("JobWorker stop")
	return j.kafkaReader.Close()

}

// indexDocument 索引文档
func (j *JobWorker) indexDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string) //这里也是去json里面找出来的
	resp, err := j.esClient.Index(j.esClient.index).
		Id(reviewID).
		Document(d).
		Do(context.Background())
	if err != nil {
		j.log.Errorf("indexing document failed, err:%v\n", err)
		return
	}
	j.log.Debugf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (j *JobWorker) updateDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	resp, err := j.esClient.Update(j.esClient.index, reviewID).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		j.log.Errorf("update document failed, err:%v\n", err)
		return
	}
	j.log.Debugf("result:%v\n", resp.Result)
}
