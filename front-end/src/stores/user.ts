import type { UserInfo } from "@/types/user";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useUserInfoStore = defineStore("userInfo", () => {
  const info = ref<UserInfo | null>(null);
  function updateInfo(data: UserInfo | null) {
    info.value = data;
  }
  return { info, updateInfo };
});
