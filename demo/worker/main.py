import os
import json
import pika
from redis.client import Redis
from pika import credentials


r: Redis = Redis(
    host=os.environ["REDIS_HOST"],
    port=int(os.environ["REDIS_PORT"]),
    db=int(os.environ["REDIS_DB"]),
    decode_responses=True
)


def job_received(channel, method, properties, body):
    try:
        job_receipt = json.loads(body)

        print("processing job id", job_receipt["job_id"])
        print("received job context", job_receipt["job_context"])

        job_key = "jobs"+job_receipt["job_id"]

        job = r.hgetall(job_key)
        if job and job.get("status") != "pending":
            return

        print(f"job {job_receipt['job_id']} still needs processing")

        message_body = {
            "job_id": job_receipt["job_id"],
            "message": "job successfully processed",
            "status": 0
        }

        channel.basic_publish(
            exchange="",
            routing_key="jobs-finished",
            body=json.dumps(message_body).encode("utf-8")
        )

        channel.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as e:
        print("error ocurred while processing job:", e)


connection = pika.BlockingConnection(pika.ConnectionParameters(
    os.environ["MQ_HOST"], credentials=credentials.PlainCredentials(os.environ["MQ_USER"], os.environ["MQ_PASS"])))

channel = connection.channel()

channel.basic_consume(
    queue=os.environ["MQ_QUEUE"],
    auto_ack=False,
    on_message_callback=job_received,
)

channel.start_consuming()
