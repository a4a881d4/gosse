typedef struct {
	volatile unsigned int slock;
} raw_spinlock_t;

static inline void __raw_spin_lock(raw_spinlock_t *lock) {
   asm("\n"
       "1:\t" "lock;" "decl %0\n"
       "jne 2f\n"
       ".align 16\n"
       "2:\trep; nop\n"
       "cmpl $0, %0\n"
       "jg 1b\n"
       "jmp 2b\n"
       ".previous"
       : "=m" (lock->slock)
       : "m" (lock->slock)
        );
}        
static inline void __raw_spin_unlock(raw_spinlock_t *lock)
{
  asm volatile("movl $1,%0" : "+m" (lock->slock) :: "memory");
}
