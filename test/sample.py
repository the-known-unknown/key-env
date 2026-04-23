import os
import time

print("\nEnv vars:")
print("TEST_CLIENT_SECRET:", os.environ.get("TEST_CLIENT_SECRET", ""))
print("TEST_CLIENT_NAME:", os.environ.get("TEST_CLIENT_NAME", ""))
print("--------------------------------")

print("\nDoing some work...")
bar_length = 10
for i in range(1, bar_length + 1):
    filled = '█' * i
    empty = '░' * (bar_length - i)
    bar = f"\r[{filled}{empty}] {i*10}%"
    print(bar, end='', flush=True)
    time.sleep(0.2)

print("\n✔ Done!\n")  # Move to next line after bar completes