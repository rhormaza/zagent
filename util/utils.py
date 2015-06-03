__doc__ = """
TODO: add docs    
"""
__author__ = "Raul Hormazabal"
__email__ = "rhormaza@gmail.com"

import os
import os.path
import fcntl

from contextlib import contextmanager

def addLocalLibPath():
    """
    Helper to add the ZenPack's lib directory to PYTHONPATH.
    """
    import os
    import site

    site.addsitedir(os.path.join(os.path.dirname(__file__), 'lib'))

def result_errmsg(result):
    """Return a useful error message string given a twisted errBack result."""
    try:

        if result.type == ConnectionRefusedError:
            return 'connection refused. Check IP and zWBEMPort'
        elif result.type == TimeoutError:
            return 'connection timeout. Check IP and zWBEMPort'
        else:
            return result.getErrorMessage()
    except AttributeError:
        pass

    return str(result)

class PersistManager(object):
    """
    This class will manage the data the we need to persist in the
    filesystem
    """

    TYPE = 'FS'
    CACHE_NAME = '/tmp/.zagent'
    LAST_QUERY_FMT = '%s/%s_%s.query'
    LAST_METADATA_FMT = '%s/%s_%s.metadata'

    @classmethod
    def create_dir(cls):
        """
        It creates the main path of our cache. Directory
        path is taken from class variable CACHE_NAME
        """
        try:
            os.makedirs(cls.CACHE_NAME)
        except OSError:
            if not os.path.isdir(cls.CACHE_NAME):
                raise

    @classmethod
    @contextmanager
    def open_and_lock(cls, filename, mode='w'):
        """
        This opens and lock a file in order to prevent
        another "thing" modifying the file at the same
        time.
        """
        f = open(filename, mode)
        try:
            # Lock the file and then yield the file reference
            fcntl.lockf(f, fcntl.LOCK_EX | fcntl.LOCK_NB)
            yield f
        except IOError:
            # We could not lock the file!
            # Then, the "finally" statement takes care
            # and we do not need to do:
            #  - fcntl.lockf(f, fcntl.LOCK_UN)
            pass
        finally:
            # Releasing and closing file
            fcntl.lockf(f, fcntl.LOCK_UN)
            f.close()

    @classmethod
    def persist_to_fs(cls, host, datasource, data, metadata=False):
        fn = None
        if metadata:
            fn = cls.LAST_METADATA_FMT % (cls.CACHE_NAME, host, datasource)
        else:
            fn = cls.LAST_QUERY_FMT % (cls.CACHE_NAME, host, datasource)
        with cls.open_and_lock(fn, 'wb') as fp:
            fp.write(data)
            #json.dump(toJson, fp, indent=2)

    @classmethod
    def read_from_fs(cls, host, datasource, metadata=False):
        fn = None
        rdata = None
        if metadata:
            fn = cls.LAST_METADATA_FMT % (cls.CACHE_NAME, host, datasource)
        else:
            fn = cls.LAST_QUERY_FMT % (cls.CACHE_NAME, host, datasource)
        #with cls.open_and_lock(filename=fn, mode='r') as fp:
        #    rdata = fp.read()
        #import pdb; pdb.set_trace()
        fp = open(fn, 'r')
        rdata = fp.read()
        fp.close()
        return rdata

if __name__ == "__main__":
     PersistManager.persist_to_fs('IP', 'DS', "{'hi':1}")
