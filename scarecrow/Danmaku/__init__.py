from abc import ABC, abstractmethod
from typing import Optional


class IDanmakuClient(ABC):
    @abstractmethod
    def start(self):
        pass

    @abstractmethod
    def stop(self):
        pass

    @abstractmethod
    def save(self, file_name: Optional[str] = None):
        pass
