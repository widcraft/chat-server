package org.wid.userservice.port.secondary;

import org.wid.userservice.entity.entity.User;

import reactor.core.publisher.Mono;

public interface UserRepositoryPort {
  Mono<User> register(User user);

  Mono<User> getUserById(String userId);
}
