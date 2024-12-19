from data.conn import Data

class SingletonMeta(type):
    _instances = {}

    def __call__(cls, *args, **kwargs):
        """
        Possible changes to the value of the `__init__` argument do not affect
        the returned instance.
        """
        if cls not in cls._instances:
            instance = super().__call__(*args, **kwargs)
            cls._instances[cls] = instance
        return cls._instances[cls]


class App(metaclass=SingletonMeta):
    def set_params(self, config, repository, project_repository):
        self.user: dict = {}
        self.config: dict = config
        self.repository: Data = repository
        self.project_repository: Data = project_repository
        self.model: dict = {}
        self.model_w: dict = {}
        self.groups: dict = {}

    def set_models(self, models):
        self.model = models
        self.make_groups()

    def make_groups(self):
        for m in self.model['models']:
            if '_to_' in m:
                l = m.split('_to_')
                if l[1] in self.groups.keys():
                    self.groups[l[1]].append(m)
                else:
                    self.groups[l[1]] = [m,]
        self.groups['project'] = ['ordering',]
