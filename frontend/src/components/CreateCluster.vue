<template>
  <div class="card">
    <h2>Create Postgres Cluster</h2>
    <form @submit.prevent="createCluster">
      <input v-model="name" placeholder="Cluster Name" required />
      <input v-model="user" placeholder="User" required />
      <input v-model="databases" placeholder="Databases (comma-separated)" required />

      <select v-model="storage" required>
        <option disabled value="">Select Storage Size</option>
        <option value="100Mi">100 MB</option>
        <option value="200Mi">200 MB</option>
        <option value="500Mi">500 MB</option>
        <option value="750Mi">750 MB</option>
        <option value="1Gi">1 Gi</option>
      </select>

      <button :disabled="loading">{{ loading ? 'Creating...' : 'Create' }}</button>
    </form>
    <p v-if="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const name = ref('')
const user = ref('')
const databases = ref('')
const storage = ref('')
const loading = ref(false)
const message = ref('')
const messageType = ref('')

const emit = defineEmits(['created'])

async function createCluster() {
  loading.value = true
  message.value = ''
  try {
    const res = await fetch('/api/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: name.value.trim(),
        user: user.value.trim(),
        databases: databases.value.split(',').map(x => x.trim()).filter(Boolean),
        storage: storage.value
      })
    })
    if (!res.ok) {
      throw new Error((await res.text()) || 'Failed to create cluster')
    }
    message.value = 'Cluster created successfully!'
    messageType.value = 'success'
    emit('created')
    name.value = user.value = databases.value = storage.value = ''
  } catch (err) {
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
  display: grid;
  gap: 0.5rem;
}
input, select {
  padding: 0.5rem;
  border-radius: 6px;
  border: 1px solid #ccc;
  outline: none;
  font-size: 1rem;
  transition: border-color 0.15s;
}
input:focus, select:focus {
  border-color: #56a1f7;
}
button {
  background: #4f8cfb;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
  margin-top: 0.5rem;
  transition: background 0.2s;
}
button:disabled {
  background: #88bffc;
  cursor: not-allowed;
}
.success {
  color: #258b40;
  margin-top: 0.3rem;
}
.error {
  color: #e25241;
  margin-top: 0.3rem;
}
</style>
