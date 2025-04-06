<script setup lang="ts">
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { UserRound, LogOut } from 'lucide-vue-next';

const user = await useUser();
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger>
      <div class="w-8 h-8 rounded-full overflow-hidden border border-muted">
        <img :src="user?.picture" :alt="user?.fullname.split(' ')[0][0]" class="w-full h-full object-cover" />
      </div>
    </DropdownMenuTrigger>
    <DropdownMenuContent>
      <DropdownMenuLabel>
        <div>
          <h1>{{ user?.fullname }}</h1>
          <p class="font-light text-muted-foreground">{{ user?.email.split('@')[0] }}</p>
        </div>
      </DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuItem><UserRound />Profile</DropdownMenuItem>
      <DropdownMenuItem><LogOut />Sign out</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>