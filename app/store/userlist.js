import { defineStore } from "pinia";
export const useListUserstore = defineStore("listusers-store", () => {
  const listUsers = ref([]);

  //actions
  const addUser = (users) => {
    listUsers.value = [...users];
  };

  const removeAllUsers = () => {
    listUsers.value = [];
  };

  return { listUsers, addUser, removeAllUsers };
});
