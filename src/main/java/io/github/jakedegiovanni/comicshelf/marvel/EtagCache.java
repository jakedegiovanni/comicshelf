package io.github.jakedegiovanni.comicshelf.marvel;

import io.github.jakedegiovanni.comicshelf.marvel.model.DataWrapper;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.locks.ReentrantReadWriteLock;

@Service
public class EtagCache<T> {

    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();
    private final Map<String, DataWrapper<T>> cache;

    public EtagCache() {
        this.cache = new HashMap<>();
    }

    public EtagCache(Map<String, DataWrapper<T>> cache) {
        this.cache = cache;
    }

    public Optional<DataWrapper<T>> get(String endpoint) {
        var rl = lock.readLock();
        rl.lock();

        try {
            return Optional.ofNullable(cache.get(endpoint));
        }
        finally {
            rl.unlock();
        }
    }

    public void put(String endpoint, DataWrapper<T> entry) {
        var wl = lock.writeLock();
        wl.lock();

        try {
            cache.put(endpoint, entry);
        }
        finally {
            wl.unlock();
        }
    }
}
