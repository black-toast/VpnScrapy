import unittest


def add(a, b):
    return a+b


class TestingMain(unittest.TestCase):
    def testAdd(self):
        test_var = add(1, 2)
        print("testvar:", test_var)
        self.assertEqual(test_var, 3)

    def testAdd2(self):
        test_var = add(2, 3)
        print("testvar2:", test_var)
        self.assertEqual(test_var, 5)


# python test_main.py TestingMain.testAdd
if __name__ == '__main__':
    unittest.main()
