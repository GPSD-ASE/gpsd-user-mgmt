from fastapi import FastAPI

app = FastAPI()

@app.get("/api/user/health")
def get_sensor_health():
    return {"status": "healthy"}
