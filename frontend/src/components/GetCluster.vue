<template>
  <div class="card">
    <h2>Get Cluster Info</h2>
    <form @submit.prevent="getCluster">
      <input v-model="name" placeholder="Cluster Name" required />
      <button :disabled="loading">{{ loading ? "Loading..." : "Get" }}</button>
    </form>
    <div v-if="infoObj" class="result">
      <div class="result-row"><strong>Name:</strong> <span>{{ infoObj.clusterName }}</span></div>
      <div class="result-row"><strong>User:</strong> <span>{{ infoObj.user }}</span></div>
      <div class="result-row"><strong>Databases:</strong> <span>{{ infoObj.databases?.join(', ') }}</span></div>
      <div class="result-row"><strong>Storage:</strong> <span>{{ infoObj.storage }}</span></div>
      <div class="result-row"><strong>NodePort:</strong> <span>{{ infoObj.nodePort || 'Pending...' }}</span></div>
    </div>
    <p v-if="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref, defineExpose } from 'vue'
const name = ref('')
const infoObj = ref(null)
const loading = ref(false)
const message = ref('')
const messageType = ref('')

defineExpose({
  showCluster(n) {
    name.value = n
    getCluster()
  }
})

async function getCluster() {
  loading.value = true
  infoObj.value = null
  message.value = ''
  try {
    const res = await fetch(`http://168.119.243.127:30002/get?name=${encodeURIComponent(name.value)}`)
    if (!res.ok) throw new Error(await res.text())
    const data = await res.json()
    infoObj.value = data
    message.value = ''
  } catch (err) {
    infoObj.value = null
    message.value = err.message
    messageType.value = 'error'
  } finally {
    loading.value = false
    setTimeout(() => (message.value = ''), 2500)
  }
}
</script>

<style scoped>
.card {
  margin-bottom: 2rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 6px #0001;
  padding: 1.5rem;
}
form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.7rem;
}
input {
  flex: 1;
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  font-size: 1rem;
}
button {
  background: #4f8cfb;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}
button:disabled {
  background: #88bffc;
  cursor: not-allowed;
}
.result {
  margin: 1rem 0 0 0;
  padding: 1rem;
  background: #f3f7fa;
  border-radius: 8px;
  font-size: 0.98rem;
}
.result-row {
  margin-bottom: 0.4rem;
  display: flex;
  gap: 1rem;
}
.result-row strong {
  width: 110px;
  display: inline-block;
}
.error {
  color: #e25241;
  margin-top: 0.3rem;
}
</style>
