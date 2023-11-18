package io.github.jakedegiovanni.comicshelf.comics;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface ComicsRepository extends JpaRepository<Comic, UUID> {
    Optional<Comic> findByInternalId(int internalId);
}
