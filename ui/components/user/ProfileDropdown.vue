<script setup lang="ts">
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { UserRound, LogOut } from 'lucide-vue-next'
import type { User } from '~/types/types'

const user = useState<User>('user')
const config = useRuntimeConfig()
const url = config.public.apiUrl + '/logout'

// fallback profile pic
const getInitial = () => user?.value?.fullname?.[0]?.toUpperCase() || '?'
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger>
      <div class="w-8 h-8 rounded-full overflow-hidden border border-border bg-primary flex items-center justify-center text-white text-sm font-medium">
        <img
          v-if="user?.picture"
          :src="user.picture"
          :alt="getInitial()"
          class="w-full h-full object-cover"
        />
        <span v-else>{{ getInitial() }}</span>
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
      <NuxtLink to="/profile">
        <DropdownMenuItem class="cursor-pointer">
          <UserRound />
          <p>Profile</p>
        </DropdownMenuItem>
      </NuxtLink>
      <NuxtLink :to="url">
        <DropdownMenuItem class="cursor-pointer">
          <LogOut />
          <p>Sign out</p>
        </DropdownMenuItem>
      </NuxtLink>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
