typedef struct {
	volatile unsigned int slock;
} raw_spinlock_t;

static inline void __raw_spin_lock(raw_spinlock_t *lock) {
   asm("\n"
       "1:\t" "lock;" "decl %0\n\t"
       "jne 2f\n\t"
       ".subsection 1\n\t"
       ".align 16\n"
       "2:\trep; nop\n\t"
       "cmpl $0, %0\n\t"
       "jg 1b\n\t"
       "jmp 2b\n\t"
       ".previous"
       : "=m" (lock->slock)
       : "m" (lock->slock)
        );
}        
static inline void __raw_spin_unlock(raw_spinlock_t *lock)
{
  asm volatile("movl $1,%0" : "+m" (lock->slock) :: "memory");
}
