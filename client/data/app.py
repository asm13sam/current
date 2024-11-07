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
    def set_params(self, config, models, repository):
        self.user: dict = {}
        self.config: dict = config
        self.model: dict = models
        self.repository: Data = repository
        self.model_w: dict = {}
        self.groups: dict = {}
        self.make_model_w()
        self.make_groups()

    def make_model_w(self):
        for m in self.model['models'].keys():
            self.model_w[m] = {}
            for k, v in self.model[m].items():
                if k.endswith('_id'):
                    table_name = '_'.join(k.split('_')[:-1])
                    if table_name.endswith('2'):
                        table_name = table_name[:-1]
                    self.model_w[m][table_name] = {}
                    if 'name' in self.model[table_name]:
                        self.model_w[m][table_name]['def'] = self.model[table_name]['name']['def']
                        self.model_w[m][table_name]['form'] = self.model[table_name]['name']['form']
                        self.model_w[m][table_name]['hum'] = v['hum']
                    
                self.model_w[m][k] = {}
                self.model_w[m][k]['def'] = v['def']
                self.model_w[m][k]['form'] = v['form']
                self.model_w[m][k]['hum'] = v['hum']

    def make_groups(self):
        for m in self.model['models']:
            if '_to_' in m:
                l = m.split('_to_')
                if l[1] in self.groups.keys():
                    self.groups[l[1]].append(m)
                else:
                    self.groups[l[1]] = [m,]
        self.groups['project'] = ['ordering',]







