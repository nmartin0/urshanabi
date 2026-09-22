"""The suite's own helpers behave as the properties expect."""

import unittest

import support


class AnswerTest(unittest.TestCase):
    """An answer reports its headers whatever their case."""

    def test_a_header_is_found_whatever_its_case(self) -> None:
        """A header asked for in any case is found."""
        answered = support.Answer(200, {"cache-control": "no-store"})
        self.assertEqual(answered.header("Cache-Control"), "no-store")

    def test_an_absent_header_is_empty_not_missing(self) -> None:
        """An absent header reads as the empty string."""
        self.assertEqual(support.Answer(404, {}).header("x-request-id"), "")


if __name__ == "__main__":
    unittest.main()
