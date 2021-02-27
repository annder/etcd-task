import http from "./config"

export const getAllTask = () => http.get("all")

export const addTask = (data) => http.post("add", data)

export const delTask = data => http.post("del", data)

export const updateStatus = data => http.post("update", data)

