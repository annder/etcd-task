<template>
  <div id="app">
    <el-button
      type="primary"
      size="default"
      @click="add"
      style="margin-bottom: 20px"
      >添加任务</el-button
    >
    <el-table :data="list" border stripe>
      <el-table-column
        v-for="col in columns"
        :prop="col.prop"
        :key="col.id"
        :label="col.label"
        align="center"
      >
      </el-table-column>
      <el-table-column label="操作">
        <template slot-scope="{ row }">
          <el-button
            v-if="row.status == 1"
            type="primary"
            icon="el-icon-video-play"
            @click="start(row)"
            >启用</el-button
          >
          <el-button
            v-else
            type="warning"
            icon="el-icon-video-pause"
            @click="puase(row)"
            >暂停</el-button
          >
          <el-button type="danger" @click="del(row)" icon="el-icon-delete-solid"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      title="添加任务"
      :visible.sync="showTask"
      width="30%"
      @close="close"
    >
      <el-form
        :model="form"
        ref="form"
        label-width="80px"
        :inline="false"
        size="normal"
      >
        <el-form-item label="任务名">
          <el-input v-model="form.name" placeholder="请输入任务名"></el-input>
        </el-form-item>

        <el-form-item label="定时器">
          <el-input
            v-model="form.cron"
            placeholder="请输入定时器"
            size="normal"
          ></el-input>
        </el-form-item>
        <el-form-item label="脚本">
          <el-input
            type="text"
            :rows="2"
            v-model="form.bash"
            placeholder="请输入脚本"
            :autosize="{ minRows: 2, maxRows: 4 }"
          >
          </el-input>
        </el-form-item>
      </el-form>

      <span slot="footer">
        <el-button @click="cancel" icon="el-icon-close">取消</el-button>
        <el-button type="primary" @click="ok" icon="el-icon-check"
          >确认</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { isValidCron } from "cron-validator";
const columns = [
  {
    label: "ID",
    prop: "ID",
  },
  {
    label: "任务名",
    prop: "name",
  },
  {
    label: "定时时间",
    prop: "cron",
  },
  {
    label: "脚本",
    prop: "bash",
  },
];
export default {
  name: "app",
  data() {
    return {
      list: [],
      form: {
        name: "",
        cron: "",
        bash: "",
      },
      showTask: false,
      columns,
    };
  },
  mounted() {
    this.getList();
  },
  methods: {
    close() {
      this.showTask = false;
    },
    async getList() {
      const { list } = await this.$api.getAllTask();
      this.list = list;
    },
    async ok() {
      if (!isValidCron(this.form.cron)) {
        this.$message.error("corn 表达式错误");
        return;
      }
      this.showTask = false;
      try {
        const data = await this.$api.addTask(this.form);
        this.list.push(data);
      } catch (e) {
        this.$message.error(e);
      }
    },
    async puase({ ID }) {
      try {
        const action = await this.$confirm("是否暂停此任务？");
        if (action === "confirm") {
          const ok = await this.$api.updateStatus({ id: ID, status: 1 });
          console.log(ok)
          if (ok) {
            this.getList();
          }
        }
      } catch (e) {
        console.log(e);
      }
    },
    // 开启
    async start({ ID }) {
      try {
        const action = await this.$confirm("是否开启此定时任务？");
        if (action === "confirm") {
          const ok = await this.$api.updateStatus({ id: ID, status: 2 }); // 2 是开启任务
          if (ok) {
            this.getList();
          }
        }
      } catch (e) {
        console.log(e);
      }
    },
    cancel() {
      this.showTask = false;
    },
    add() {
      this.showTask = true;
    },
    /**
     * @name 删除一个任务
     * */

    async del({ ID }) {
      try {
        const action = await this.$confirm("确认删除此任务？");
        if (action === "confirm") {
          const data = await this.$api.delTask({ id: ID });
          if (data) {
            this.getList();
          }
        }
      } catch (e) {
        this.$message.info("取消删除");
      }
    },
  },
};
</script>
