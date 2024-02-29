from abc import ABC, abstractmethod

class RepoFactory(ABC):
    @staticmethod
    @abstractmethod
    def get_scoped_conn(url):
        pass
    
    @classmethod
    @abstractmethod
    def get_singleton_conn(self, url):
        pass