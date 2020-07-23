"""用户自定义的orm映射.

注意使用装饰器register将要创建的表注册到Tables.
"""
from io import BytesIO
import pyotp
from qrcode import QRCode, constants
from peewee import (
    CharField,
    IntegerField,
    AutoField,
    BooleanField
)

from const import SERVICE_NAME
from ._base import (
    BaseModel,
    register
)


@register
class User(BaseModel):
    """用户的数据模型."""
    uid = AutoField()
    name = CharField()
    email = CharField()
    password = CharField()
    is_admin = BooleanField(default=False)
    gtoken = CharField(default="")

    def secondary_verify(self,code:int)->bool:
        """两步验证.

        Args:
            code (int): google 验证器的动态校验码

        Returns:
            bool: 是否符合
        """
        t = pyotp.TOTP(self.gtoken)
        result = t.verify(code)
        msg = result if result is True else False
        return msg

    def new_google_auth_qrcode(self)->str:
        """创建googleauth使用的二维码.

        Returns:
            str: 二维码的图片base64字符串
        """
        data = pyotp.totp.TOTP(self.gtoken).provisioning_uri(self.name, issuer_name=SERVICE_NAME)
        qr = QRCode(
        version=1,
        error_correction=constants.ERROR_CORRECT_L,
        box_size=6,
        border=4, )

        qr.add_data(data)
        qr.make(fit=True)
        img = qr.make_image()
        output_buffer = BytesIO()
        img.save(output_buffer, format='png')
        byte_data = output_buffer.getvalue()
        base64_str = base64.b64encode(byte_data)
        return f"data:image/png;base64,{base64_str}"

    def new_secondary_verification_token(self):
        if self.gtoken == "":
            self.gtoken =  pyotp.random_base32(64)
        else:
            raise 


    def to_dict(self):
        return {
            "uid": self.uid,
            "name": self.name,
            "email": self.email,
            "is_admin": self.is_admin
        }


__all__ = ["User"]
