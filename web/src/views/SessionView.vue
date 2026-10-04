<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { addSession, removeSession, selectSession, sessionStore } from '@/sessions'

const token = ref('')
const alias = ref('')

const mask = (t: string) => `${t.slice(0, 3)}***${t.slice(-3)}`

function add() {
  if (!token.value.trim()) {
    ElMessage.warning('请输入 Session')
    return
  }
  addSession(token.value.trim(), alias.value.trim())
  token.value = ''
  alias.value = ''
  ElMessage.success('已添加')
}

function copy(t: string) {
  navigator.clipboard.writeText(t)
  ElMessage.success('Session 已复制到剪贴板')
}

async function del(index: number) {
  await ElMessageBox.confirm('确定要删除这个 Session 吗？', '提示', { type: 'warning' })
  removeSession(index)
  ElMessage.success('已删除')
}
</script>

<template>
  <div>
    <h1 class="mb-4 text-2xl font-bold text-gray-200">Sessions</h1>

    <div class="frost mb-4 flex flex-wrap items-center gap-2 p-3 sm:p-4">
      <el-input v-model="token" placeholder="输入一个 Session..." class="w-full sm:w-72" />
      <el-input v-model="alias" placeholder="输入别名..." class="w-full sm:w-48" />
      <el-button type="primary" @click="add">添加 Session</el-button>
    </div>

    <el-empty v-if="sessionStore.list.length === 0" description="还没有添加任何 Session 哦" />

    <div v-else class="frost p-2 sm:p-3">
      <el-table :data="sessionStore.list" stripe>
        <el-table-column label="选择" width="90">
          <template #default="{ row }">
            <el-button
              size="small"
              :type="sessionStore.selected === row.token ? 'primary' : 'default'"
              @click="selectSession(row.token)"
            >
              {{ sessionStore.selected === row.token ? '已选' : '选择' }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="Session">
          <template #default="{ row }">{{ mask(row.token) }}</template>
        </el-table-column>
        <el-table-column prop="alias" label="别名" />
        <el-table-column label="操作" width="220">
          <template #default="{ row, $index }">
            <el-button size="small" @click="copy(row.token)">复制</el-button>
            <el-button size="small" type="danger" @click="del($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>