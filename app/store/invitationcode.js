import { defineStore } from "pinia";
export const useInvitationCodeStore = defineStore(
  "invitationCode-store",
  () => {
    const invitationCode = ref();
    return { invitationCode };
  },

  {
    persist: true,
  }
);
