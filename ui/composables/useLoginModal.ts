import { ref } from 'vue'

const showLoginModal = ref(false)

export function useLoginModal() {
  return {
    showLoginModal,
    openLoginModal: () => (showLoginModal.value = true),
    closeLoginModal: () => (showLoginModal.value = false),
  }
}